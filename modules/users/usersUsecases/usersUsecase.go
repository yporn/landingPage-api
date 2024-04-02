package usersUsecases

import (
	"fmt"
	"math"
	"strconv"

	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/users"
	"github.com/yporn/sirarom-backend/modules/users/usersRepositories"
	"github.com/yporn/sirarom-backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	FindOneUser(userId string) (*users.User, error)
	FindUser(req *users.UserFilter) *entities.PaginateRes
	InsertAdmin(req *users.User) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	UpdateUser(req *users.User) (*users.User, error)
	DeleteUser(userId string) error
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) FindOneUser(userId string) (*users.User, error) {
	user, err := u.usersRepository.FindOneUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usersUsecase) FindUser(req *users.UserFilter) *entities.PaginateRes {
	users, count := u.usersRepository.FindUser(req)

	return &entities.PaginateRes{
		Data:      users,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}


func (u *usersUsecase) InsertAdmin(req *users.User) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.usersRepository.InsertUser(req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	//Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Compare Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	// Sign token
	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:       user.Id,
		UserRole: make([]*users.UserRole, 0),
	})
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshToken, err := auth.NewAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:       user.Id,
		UserRole: make([]*users.UserRole, 0),
	})
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(user.Id)
	// Set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       id,
			Email:    user.Email,
			Username: user.Username,
			Images:   make([]*entities.Image, 0),
			UserRole: make([]*users.UserRole, 0),
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Check Oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:       strconv.Itoa(profile.Id),
		UserRole: make([]*users.UserRole, 0),
	}

	accessToken, err := auth.NewAuth(
		auth.Access,
		u.cfg.Jwt(),
		newClaims,
	)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}

	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}

func (u *usersUsecase) UpdateUser(req *users.User) (*users.User, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	user, err := u.usersRepository.UpdateUser(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usersUsecase) DeleteUser(userId string) error {
	if err := u.usersRepository.DeleteUser(userId); err != nil {
		return err
	}
	return nil
}
