package jobs

import "github.com/yporn/sirarom-backend/modules/entities"

type Job struct {
	Id            int `db:"id" json:"id"`
	Position      string `db:"position" json:"position"`
	Amount        string `db:"amount" json:"amount"`
	Location      string `db:"location" json:"location"`
	Description   string `db:"description" json:"description"`
	Qualification string `db:"qualification" json:"qualification"`
	StartDate     string `db:"start_date" json:"start_date"`
	EndDate       string `db:"end_date" json:"end_date"`
	Status        string `db:"status" json:"status"`
	Display       string `db:"display" json:"display"`
	CreatedAt     string `db:"created_at" json:"created_at"`
	UpdatedAt     string `db:"updated_at" json:"updated_at"`
}

type JobFilter struct {
	Id string `query:"id"`
	Search string `query:"search"`
	*entities.PaginationReq
	*entities.SortReq
}