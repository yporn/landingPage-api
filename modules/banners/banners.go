package banners

import "github.com/yporn/sirarom-backend/modules/entities"

type Banner struct {
	Id      int               `db:"id" json:"id"`
	Index   int               `db:"index" json:"index"`
	Delay   int               `db:"delay" json:"delay"`
	Display string            `db:"display" json:"display"`
	Images  []*entities.Image `json:"images"`
}

type BannerFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // title & description
	*entities.PaginationReq
	*entities.SortReq
}