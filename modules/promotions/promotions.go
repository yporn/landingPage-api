package promotions

import (
	"github.com/yporn/sirarom-backend/modules/entities"
)

type Promotion struct {
	Id          int                    `db:"id" json:"id" form:"id"`
	Index       int                    `db:"index" json:"index" form:"index" `
	Heading     string                 `db:"heading" json:"heading" form:"heading"`
	Description string                 `db:"description" json:"description" form:"description"`
	StartDate   string                 `db:"start_date" json:"start_date" form:"start_date"`
	EndDate     string                 `db:"end_date" json:"end_date" form:"end_date"`
	Display     string                 `db:"display" json:"display" form:"display"`
	Images      []*entities.Image      `json:"promotion_images"`
	HouseModel  []*PromotionHouseModel `json:"house_models"`
	FreeItem    []*PromotionFreeItem   `json:"free_items"`
}

type PromotionFreeItem struct {
	Id          int    `db:"id" json:"id"`
	PromotionId int    `db:"promotion_id" json:"promotion_id"`
	Description string `db:"description" json:"description"`
}

type PromotionHouseModel struct {
	Id           int           `db:"id" json:"id"`
	PromotionId  int           `db:"promotion_id" json:"promotion_id"`
	HouseModelId int           `db:"house_model_id" json:"house_model_id"`
	HouseModel   []*HouseModel `json:"house_model_name"`
}

type HouseModel struct {
	Id        int               `json:"id"`
	ProjectId int               `json:"project_id"`
	Name      string            `json:"name"`
	HouseType []*HouseModelType `json:"house_type"`
	Images    []*entities.Image `json:"house_images"`
}

type HouseModelType struct {
	Id       int    `json:"id"`
	RoomType string `json:"room_type"`
	Amount   int    `json:"amount"`
}

type PromotionFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"` // Heading
	*entities.PaginationReq
	*entities.SortReq
}
