package promotions

import (
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/houseModels"
)

type Promotion struct {
	Id          int               `db:"id" json:"id"`
	Index       int               `db:"index" json:"index"`
	Heading     string            `db:"heading" json:"heading"`
	Description string            `db:"description" json:"description"`
	StartDate   string            `db:"start_date" json:"start_date"`
	EndDate     string            `db:"end_date" json:"end_date"`
	Display     string            `db:"display" json:"display"`
	Images      []*entities.Image `json:"house_images"`
	HouseModel  []HouseModel      `json:"house_model"`
}

type HouseModel struct {
	Id         int                     `db:"id" json:"id"`
	HouseModel *houseModels.HouseModel `db:"house_model" json:"house_model"`
}

type PromotionFreeItem struct {
	Id          int    `db:"id" json:"id"`
	PromotionId int    `db:"promotion_id" json:"promotion_id"`
	Description string `db:"description" json:"description"`
}
