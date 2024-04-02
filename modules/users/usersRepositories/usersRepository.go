package usersRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/users"
	"github.com/yporn/sirarom-backend/modules/users/usersPatterns"
)

type IUsersRepository interface {
	InsertUser(req *users.User) (*users.UserPassport, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	FindUser(req *users.UserFilter) ([]*users.User, int)
	InsertOauth(req *users.UserPassport) error
	FindOneOauth(refreshToken string) (*users.Oauth, error)
	FindOneUser(userId string) (*users.User, error)
	UpdateOauth(req *users.UserToken) error
	GetProfile(userId string) (*users.User, error)
	DeleteOauth(oauthId string) error
	UpdateUser(req *users.User) (*users.User, error)
	DeleteUser(userId string) error
}

type usersRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesUsecases.IFilesUsecase
}

func UsersRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesUsecases.IFilesUsecase) IUsersRepository {
	return &usersRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *usersRepository) InsertUser(req *users.User) (*users.UserPassport, error) {
	builder := usersPatterns.InsertUserBuilder(r.db, req)

	// Insert the user into the database
	if _, err := usersPatterns.InsertUserEngineer(builder).InsertUser(); err != nil {
		return nil, fmt.Errorf("insert user failed: %v", err)
	}

	// Retrieve the inserted user data
	user, err := builder.Result()
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return nil, fmt.Errorf("inserted user not found: %v", err)
		}

		return nil, fmt.Errorf("get user failed: %v", err)
	}

	return user, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"id",
		"email",
		"password",
		"username"
	FROM "users"
	WHERE "email" = $1;`

	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, email); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *usersRepository) FindOneUser(userId string) (*users.User, error) {
	query := `
	SELECT to_jsonb("t")
	FROM (
		SELECT
			"u".*,
			(
				SELECT
					COALESCE(array_to_json(array_agg("i")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "user_images" "i"
					WHERE "i"."user_id" = "u"."id"
				) AS "i"
			) AS "images",
			(
				SELECT
					COALESCE(array_to_json(array_agg("r")), '[]'::json)
				FROM (
					SELECT
						"r"."id",
						"r"."role_id"
					FROM "user_roles" "r"
					WHERE "r"."user_id" = "u"."id"
				) AS "r"
			) AS "roles"
		FROM "users" "u"
		WHERE "u"."id" = $1
	) AS "t";`

	userBytes := make([]byte, 0)
	user := &users.User{
		Images:   make([]*entities.Image, 0),
		UserRole: make([]*users.UserRole, 0),
	}

	if err := r.db.Get(&userBytes, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}
	return user, nil
}

func (r *usersRepository) FindUser(req *users.UserFilter) ([]*users.User, int) {
	builder := usersPatterns.FindUserBuilder(r.db, req)
	engineer := usersPatterns.FindUserEngineer(builder)

	result := engineer.FindUser().Result()
	count := engineer.CountUser().Count()
	return result, count
}

func (r *usersRepository) InsertOauth(req *users.UserPassport) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := `
	INSERT INTO "oauth" (
		"user_id",
		"refresh_token",
		"access_token"
	)
	VALUES ($1, $2, $3)
		RETURNING "id";
	`
	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User.Id,
		req.Token.RefreshToken,
		req.Token.AccessToken,
	).Scan(&req.Token.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) FindOneOauth(refreshToken string) (*users.Oauth, error) {
	query := `
	SELECT
		"id",
		"user_id"
	FROM "oauth"
	WHERE "refresh_token" = $1;
	`
	oauth := new(users.Oauth)
	if err := r.db.Get(oauth, query, refreshToken); err != nil {
		return nil, fmt.Errorf("oauth not found")
	}
	return oauth, nil
}

func (r *usersRepository) UpdateOauth(req *users.UserToken) error {
	query := `
	UPDATE "oauth" SET
		"access_token" = :access_token,
		"refresh_token" = :refresh_token
	WHERE "id" = :id;
	`

	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}

	return nil
}

func (r *usersRepository) GetProfile(userId string) (*users.User, error) {
	query := `
	SELECT
		"id",
		"email",
		"username"
		
		
	FROM "users"
	WHERE "id" = $1;`

	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}
	return profile, nil
}

func (r *usersRepository) DeleteOauth(oauthId string) error {
	query := `DELETE FROM "oauth" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, oauthId); err != nil {
		return fmt.Errorf("oauth not found: %v", err)
	}
	return nil
}

func (r *usersRepository) UpdateUser(req *users.User) (*users.User, error) {
	builder := usersPatterns.UpdateUserBuilder(r.db, req, r.filesUsecase)
	engineer := usersPatterns.UpdateUserEngineer(builder)

	if err := engineer.UpdateUser(); err != nil {
		return nil, err
	}

	user, err := r.FindOneUser(strconv.Itoa(req.Id))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) DeleteUser(userId string) error {
	query := `DELETE FROM "users" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, userId); err != nil {
		return fmt.Errorf("delete users failed: %v", err)
	}
	return nil
}
