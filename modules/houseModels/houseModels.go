package houseModels

import "github.com/yporn/sirarom-backend/modules/entities"

type HouseModel struct {
	Id              int                   `db:"id" json:"id"`
	ProjectId       int                   `db:"project_id" json:"project_id"`
	Name            string                `db:"name" json:"name"`
	Description     string                `db:"description"`
	LinkVideo       string                `db:"link_video" json:"link_video"`
	LinkVirtualTour string                `db:"link_virtual_tour" json:"link_virtual_tour"`
	Display         string                `db:"display" json:"display"`
	Index           int                   `db:"index" json:"index"`
	CreatedAt       string                `db:"created_at" json:"created_at"`
	UpdatedAt       string                `db:"updated_at" json:"updated_at"`
	Images          []*entities.Image     `json:"house_images"`
	TypeItem        []*HouseModelTypeItem `json:"type_items"`
	HousePlan       []*HouseModelPlan     `json:"house_plan"`
}

type HouseModelTypeItem struct {
	Id       int    `db:"id" json:"id"`
	RoomType string `db:"room_type" json:"room_type"`
	Amount   int `db:"amount" json:"amount"`
}

type HouseModelPlan struct {
	Id           int                   `db:"id" json:"id"`
	HouseModelId int                   `db:"house_model_id" json:"house_model_id"`
	Floor        int                `db:"floor" json:"floor"`
	Size         string                `db:"size" json:"size"`
	Images       []*entities.Image     `json:"plan_images"`
	PlanItem     []*HouseModelPlanItem `json:"plan_items"`
}

type HouseModelPlanItem struct {
	Id          int    `db:"id" json:"id"`
	HousePlanId int    `db:"house_model_plan_id" json:"house_model_plan_id"`
	RoomType    string `db:"room_type" json:"room_type"`
	Amount      int `db:"amount" json:"amount"`
}

type HouseModelFilter struct {
	Id        string `query:"id"`
	ProjectId int    `query:"project_id"`
	Search    string `query:"search"` // name
	*entities.PaginationReq
	*entities.SortReq
}

type HouseModelName struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
}