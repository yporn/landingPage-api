package projects

import (
	"encoding/json"

	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/houseModels"
)

type Project struct {
	Id            int                       `db:"id" json:"id"`
	Name          string                    `db:"name" json:"name"`
	Index         int                       `db:"index" json:"index"`
	Heading       string                    `db:"heading" json:"heading"`
	Text          string                    `db:"text" json:"text"`
	Location      string                    `db:"location" json:"location"`
	Price         int                       `db:"price" json:"price"`
	StatusProject string                    `db:"status_project" json:"status_project"`
	TypeProject   string                    `db:"type_project" json:"type_project"`
	Description   string                    `db:"description" json:"description"`
	NameFacebook  string                    `db:"name_facebook" json:"name_facebook"`
	LinkFacebook  string                    `db:"link_facebook" json:"link_facebook"`
	Tel           string                    `db:"tel" json:"tel"`
	Address       string                    `db:"address" json:"address"`
	LinkLocation  string                    `db:"link_location" json:"link_location"`
	Display       string                    `db:"display" json:"display"`
	CreatedAt     string                    `db:"created_at" json:"created_at"`
	UpdatedAt     string                    `db:"updated_at" json:"updated_at"`
	Images        []*entities.Image         `json:"images"`
	HouseTypeItem []*ProjectHouseTypeItem   `json:"house_type_items"`
	DescAreaItem  []*ProjectDescAreaItem    `json:"area_items"`
	FacilityItem  []*ProjectFacilityItem    `json:"facilities_items"`
	HouseModel    []*houseModels.HouseModel `json:"house_models"`
}

type ProjectHouseTypeItem struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type ProjectDescAreaItem struct {
	Id       int    `db:"id" json:"id"`
	ItemArea string `db:"item" json:"item"`
	Amount   int    `db:"amount" json:"amount"`
	Unit     string `db:"unit" json:"unit"`
}

// facilities
type ProjectFacilityItem struct {
	Id   int    `db:"id" json:"id"`
	Item string `db:"item" json:"item"`
}

type ProjectFilter struct {
	Search        string `query:"search"` // name,status_project,type_project,location
	StatusProject string `query:"status_project"`
	*entities.PaginationReq
	*entities.SortReq
}
type ProjectHouseModelResult struct {
	ProjectHouseModel json.RawMessage `db:"project_house_model"`
}
