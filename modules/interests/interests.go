package interests

import "github.com/yporn/sirarom-backend/modules/entities"

type Interest struct {
	Id           int    `db:"id" json:"id"`
	BankName     string `db:"bank_name" json:"bank_name"`
	InterestRate string `db:"interest_rate" json:"interest_rate"`
	Note         string `db:"note" json:"note"`
	Display      string `db:"display" json:"display"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	UpdatedAt    string `db:"updated_at" json:"updated_at"`
	// FileName     string `db:"filename" json:"filename"`
	// Url          string `db:"url" json:"url"`
	Images      []*entities.Image `json:"images"`
}

type InterestFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // title & description
	*entities.PaginationReq
	*entities.SortReq
}
