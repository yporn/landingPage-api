package users

import (
	"fmt"
	"regexp"

	"github.com/yporn/sirarom-backend/modules/entities"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int               `db:"id" json:"id"`
	Email    string            `db:"email" json:"email"`
	Username string            `db:"username" json:"username"`
	Password string            `db:"password" json:"password"`
	Name     string            `db:"name" json:"name"`
	Tel      string            `db:"tel" json:"tel"`
	Display  string            `db:"display" json:"display"`
	Images   []*entities.Image `json:"images"`
	UserRole []*UserRole       `json:"roles"`
}

type UserRegisterReq struct {
	Email    string            `db:"email" json:"email" form:"email"`
	Password string            `db:"password" json:"password" form:"password"`
	Username string            `db:"username" json:"username" form:"username"`
	Name     string            `db:"name" json:"name"`
	Tel      string            `db:"tel" json:"tel"`
	Display  string            `db:"display" json:"display"`
	Images   []*entities.Image `json:"images"`
	UserRole []*UserRole       `json:"roles"`
}

type UserRole struct {
	Id     int `db:"id" json:"id"`
	UserId int `db:"user_id" json:"user_id"`
	RoleId int `db:"role_id" json:"role_id"`
}

type UserCredential struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserCredentialCheck struct {
	Id       string      `db:"id"`
	Email    string      `db:"email"`
	Password string      `db:"password"`
	Username string      `db:"username"`
	UserRole []*UserRole `json:"roles"`
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `db:"id" json:"id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

type UserClaims struct {
	Id       string      `db:"id" json:"id"`
	UserRole []*UserRole `json:"roles"`
}

type UserRefreshCredential struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type Oauth struct {
	Id     string `db:"id" json:"id"`
	UserId string `db:"user_id" json:"user_id"`
}

type UserRemoveCredential struct {
	OauthId string `db:"id" json:"oauth_id" form:"oauth_id"`
}

func (obj *User) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hashed password failed: %v", err)
	}
	obj.Password = string(hashedPassword)
	return nil
}

func (obj *User) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}

type UserFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // name & username & email
	*entities.PaginationReq
	*entities.SortReq
}