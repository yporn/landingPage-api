package middlewaresRepositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yporn/sirarom-backend/modules/middlewares"
)

type IMiddlewaresRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
	GetUserRoles(userID int) ([]*middlewares.Role, error) 
}

type middlewaresRepository struct {
	db *sqlx.DB
}

func MiddlewaresRepository(db *sqlx.DB) IMiddlewaresRepository {
	return &middlewaresRepository{
		db: db,
	}
}

func (r *middlewaresRepository) FindAccessToken(userId, accessToken string) bool {
	query := `
	SELECT
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
	FROM "oauth"
	WHERE "user_id" = $1
	AND "access_token" = $2;
	`
	var check bool
	if err := r.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}
	return true
}

func (r *middlewaresRepository) FindRole() ([]*middlewares.Role, error) {
	query := `
	SELECT
		"id",
		"title"
	FROM "roles"
	ORDER BY "id" DESC;`

	roles := make([]*middlewares.Role, 0)
	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empty")
	}
	return roles, nil
}


func (r *middlewaresRepository) GetUserRoles(userID int) ([]*middlewares.Role, error) {
	query := `
	SELECT
		r.id,
		r.title,
		ur.id AS user_role_id,
		ur.user_id,
		ur.role_id
	FROM roles r
	LEFT JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1
	ORDER BY r.id DESC;`

rolesMap := make(map[int]*middlewares.Role)
var roles []*middlewares.Role

rows, err := r.db.Query(query, userID)
if err != nil {
	return nil, err
}
defer rows.Close()

for rows.Next() {
	var roleId int
	var title string
	var userRoleID int
	var userID int
	var roleID int
	if err := rows.Scan(&roleId, &title, &userRoleID, &userID, &roleID); err != nil {
		return nil, err
	}
	role, ok := rolesMap[roleID]
	if !ok {
		role = &middlewares.Role{
			Id:    roleID,
			Title: title,
		}
		rolesMap[roleID] = role
		roles = append(roles, role)
	}
	role.UserRoles = append(role.UserRoles, &middlewares.UserRole{
		Id:     userRoleID,
		UserId: userID,
		RoleId: roleID,
	})
}

if err := rows.Err(); err != nil {
	return nil, err
}

return roles, nil
}