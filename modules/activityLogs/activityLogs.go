package activityLogs

type ActivityLog struct {
	Id     int    `json:"id"`
	User   *User  `json:"user"`
	Action string `json:"action"`
	Detail string `json:"details"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
}
