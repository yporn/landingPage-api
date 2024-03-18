package houseModelsHandlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/houseModels"
	"github.com/yporn/sirarom-backend/modules/houseModels/houseModelsUsecases"
)

type houseModelsHandlersErrCode string

const (
	findOneHouseModelErr houseModelsHandlersErrCode = "houses-001"
	findHouseModelErr    houseModelsHandlersErrCode = "houses-002"
	insertHouseModelErr  houseModelsHandlersErrCode = "houses-003"
	deleteHouseModelErr  houseModelsHandlersErrCode = "houses-004"
	updateHouseModelErr  houseModelsHandlersErrCode = "houses-005"
)

type IHouseModelsHandler interface {
	FindOneHouseModel(c *fiber.Ctx) error
	FindHouseModel(c *fiber.Ctx) error 
	AddHouseModel(c *fiber.Ctx) error
}

type houseModelsHandler struct {
	cfg                 config.IConfig
	houseModelsUsecases houseModelsUsecases.IHouseModelsUsecase
	filesUsecase        filesUsecases.IFilesUsecase
}

func HouseModelsHandler(cfg config.IConfig, houseModelsUsecase houseModelsUsecases.IHouseModelsUsecase, filesUsecase filesUsecases.IFilesUsecase) IHouseModelsHandler {
	return &houseModelsHandler{
		cfg:                 cfg,
		houseModelsUsecases: houseModelsUsecase,
		filesUsecase:        filesUsecase,
	}
}

func (h *houseModelsHandler) FindOneHouseModel(c *fiber.Ctx) error {
	houseId := strings.Trim(c.Params("house_model_id"), " ")

	house, err := h.houseModelsUsecases.FindOneHouseModel(houseId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneHouseModelErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, house).Res()
}

func (h *houseModelsHandler) FindHouseModel(c *fiber.Ctx) error {
    projectIdStr := c.Params("project_id")

    projectId, err := strconv.Atoi(projectIdStr)
    if err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            "Invalid project ID",
            err.Error(),
        ).Res()
    }

    // Create a pointer to HouseModelFilter struct
    req := &houseModels.HouseModelFilter{
        SortReq:       &entities.SortReq{},
        PaginationReq: &entities.PaginationReq{},
    }

    // Parse query parameters into the req struct (which is a pointer)
    if err := c.QueryParser(req); err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(findHouseModelErr),
            err.Error(),
        ).Res()
    }

    // Set the project ID in the filter
    req.ProjectId = projectId

    // Paginate
    if req.Page < 1 {
        req.Page = 1
    }
    if req.Limit < 5 {
        req.Limit = 100000000
    }

    // Sort
    orderByMap := map[string]string{
        "id":         `"hm"."id"`,
        "created_at": `"hm"."created_at"`,
    }
    if orderByMap[req.OrderBy] == "" {
        req.OrderBy = orderByMap["id"]
    }

    req.Sort = strings.ToUpper(req.Sort)
    sortMap := map[string]string{
        "DESC": "DESC",
        "ASC":  "ASC",
    }
    if sortMap[req.Sort] == "" {
        req.Sort = sortMap["DESC"]
    }

    // Retrieve house models for the specified project ID and filter
	houseModels := h.houseModelsUsecases.FindHouseModel(projectIdStr, req)
    // Return response
    return entities.NewResponse(c).Success(
        fiber.StatusOK,
        houseModels,
    ).Res()
}


func (h *houseModelsHandler) AddHouseModel(c *fiber.Ctx) error {
	req := &houseModels.HouseModel{
		TypeItem:  make([]*houseModels.HouseModelTypeItem, 0),
		HousePlan: make([]*houseModels.HouseModelPlan, 0),
		Images:    make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertHouseModelErr),
			err.Error(),
		).Res()
	}

	if len(req.TypeItem) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertHouseModelErr),
			"house type item id is invalid",
		).Res()
	}

	if len(req.HousePlan) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertHouseModelErr),
			"area item id is invalid",
		).Res()
	}

	houseModel, err := h.houseModelsUsecases.AddHouseModel(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertHouseModelErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, houseModel).Res()
}
