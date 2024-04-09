package middlewares

type Role struct {
	Id    int    `db:"id"`
	Title string `db:"title"`
	UserRoles []*UserRole
}

type UserRole struct {
	Id     int `db:"id" json:"id"`
	UserId int `db:"user_id" json:"user_id"`
	RoleId int `db:"role_id" json:"role_id"`
}