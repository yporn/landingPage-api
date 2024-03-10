package activities

import "github.com/yporn/sirarom-backend/modules/entities"

type Activity struct {
	Id          int               `db:"id" json:"id"`
	Index       int               `db:"index" json:"index"`
	Heading     string            `db:"heading" json:"heading"`
	Description string            `db:"description" json:"description"`
	StartDate   string            `db:"start_date" json:"start_date"`
	EndDate     string            `db:"end_date" json:"end_date"`
	VideoLink   string            `db:"video_link" json:"video_link"`
	Display     string            `db:"display" json:"display"`
	Images      []*entities.Image `json:"images"`
}

type ActivityFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // title & description
	*entities.PaginationReq
	*entities.SortReq
}
