package logos

import "github.com/yporn/sirarom-backend/modules/entities"

type Logo struct {
	Id      int               `db:"id" json:"id"`
	Index   int               `db:"index" json:"index"`
	Name    string            `db:"name" json:"name"`
	Display string            `db:"display" json:"display"`
	Images  []*entities.Image `json:"images"`
}

type LogoFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // name
	*entities.PaginationReq
	*entities.SortReq
}
