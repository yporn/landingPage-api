package houseModelsHandlers

import (
	"database/sql"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/houseModels"
	"github.com/yporn/sirarom-backend/modules/houseModels/houseModelsUsecases"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type houseModelsHandlersErrCode string

const (
	findOneHouseModelErr houseModelsHandlersErrCode = "houses-001"
	findHouseModelErr    houseModelsHandlersErrCode = "houses-002"
	findAllHouseModelErr houseModelsHandlersErrCode = "houses-006"
	insertHouseModelErr  houseModelsHandlersErrCode = "houses-003"
	deleteHouseModelErr  houseModelsHandlersErrCode = "houses-004"
	updateHouseModelErr  houseModelsHandlersErrCode = "houses-005"
)

type IHouseModelsHandler interface {
	FindOneHouseModel(c *fiber.Ctx) error
	FindHouseModel(c *fiber.Ctx) error
	FindAllHouseModel(c *fiber.Ctx) error
	AddHouseModel(c *fiber.Ctx) error
	UpdateHouseModel(c *fiber.Ctx) error
	DeleteHouseModel(c *fiber.Ctx) error
}

type houseModelsHandler struct {
	cfg                 config.IConfig
	houseModelsUsecases houseModelsUsecases.IHouseModelsUsecase
	filesUsecase        filesUsecases.IFilesUsecase
	db                  *sql.DB
}

func HouseModelsHandler(cfg config.IConfig, houseModelsUsecase houseModelsUsecases.IHouseModelsUsecase, filesUsecase filesUsecases.IFilesUsecase, db *sql.DB) IHouseModelsHandler {
	return &houseModelsHandler{
		cfg:                 cfg,
		houseModelsUsecases: houseModelsUsecase,
		filesUsecase:        filesUsecase,
		db:                  db,
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

func (h *houseModelsHandler) FindAllHouseModel(c *fiber.Ctx) error {
	houses, err := h.houseModelsUsecases.FindAllHouseModels()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findAllHouseModelErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, houses).Res()
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

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "created", "เพิ่มข้อมูลแบบบ้าน : "+houseModel.Name)
	if err != nil {
		// Handle error if logging fails
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			fmt.Sprintf("Failed to log activity %v", userID),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, houseModel).Res()
}

func (h *houseModelsHandler) UpdateHouseModel(c *fiber.Ctx) error {
	houseModelIdStr := strings.Trim(c.Params("house_model_id"), " ")
	houseModelId, err := strconv.Atoi(houseModelIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateHouseModelErr),
			err.Error(),
		).Res()
	}

	req := &houseModels.HouseModel{
		TypeItem:  make([]*houseModels.HouseModelTypeItem, 0),
		HousePlan: make([]*houseModels.HouseModelPlan, 0),
		Images:    make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateHouseModelErr),
			err.Error(),
		).Res()
	}
	req.Id = houseModelId

	houseModel, err := h.houseModelsUsecases.UpdateHouseModel(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateHouseModelErr),
			err.Error(),
		).Res()
	}

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "updated", "แก้ไขข้อมูลแบบบ้าน : "+houseModel.Name)
	if err != nil {
		// Handle error if logging fails
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			fmt.Sprintf("Failed to log activity %v", userID),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, houseModel).Res()
}

func (h *houseModelsHandler) DeleteHouseModel(c *fiber.Ctx) error {
	houseId := strings.Trim(c.Params("house_model_id"), " ")

	houseModel, err := h.houseModelsUsecases.FindOneHouseModel(houseId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteHouseModelErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range houseModel.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("house_models/%s", path.Base(p.Url)),
		})
	}

	for _, housePlan := range houseModel.HousePlan {
		for _, image := range housePlan.Images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("house_model_plans/%s", path.Base(image.Url)),
			})
		}
	}

	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteHouseModelErr),
			err.Error(),
		).Res()
	}

	if err := h.houseModelsUsecases.DeleteHouseModel(houseId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteHouseModelErr),
			err.Error(),
		).Res()
	}

		// Log activity
		userID := utils.GetUserIDFromContext(c)
		err = utils.LogActivity(h.db, strconv.Itoa(userID), "deleted", "ลบข้อมูลแบบบ้าน : "+houseModel.Name)
		if err != nil {
			// Handle error if logging fails
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				fmt.Sprintf("Failed to log activity %v", userID),
				err.Error(),
			).Res()
		}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
