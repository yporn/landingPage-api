package promotionsHandlers

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/promotions"
	"github.com/yporn/sirarom-backend/modules/promotions/promotionsUsecases"
)

type promotionsHandlersErrCode string

const (
	findOnePromotionErr promotionsHandlersErrCode = "promotions-001"
	findPromotionErr    promotionsHandlersErrCode = "promotions-002"
	insertPromotionErr  promotionsHandlersErrCode = "promotions-003"
	deletePromotionErr  promotionsHandlersErrCode = "promotions-004"
	updatePromotionErr  promotionsHandlersErrCode = "promotions-005"
)

type IPromotionsHandler interface {
	FindOnePromotion(c *fiber.Ctx) error
	FindPromotion(c *fiber.Ctx) error 
	AddPromotion(c *fiber.Ctx) error
	UpdatePromotion(c *fiber.Ctx) error
	DeletePromotion(c *fiber.Ctx) error 
}

type promotionsHandlers struct {
	cfg               config.IConfig
	promotionsUsecase promotionsUsecases.IPromotionsUsecase
	filesUsecase      filesUsecases.IFilesUsecase
}

func PromotionsHandler(cfg config.IConfig, promotionsUsecase promotionsUsecases.IPromotionsUsecase, filesUsecase filesUsecases.IFilesUsecase) IPromotionsHandler {
	return &promotionsHandlers{
		cfg: cfg,
		promotionsUsecase: promotionsUsecase,
		filesUsecase: filesUsecase,
	}
}

func (h *promotionsHandlers) FindOnePromotion(c *fiber.Ctx) error {
	promotionId := strings.Trim(c.Params("promotion_id"), " ")

	house, err := h.promotionsUsecase.FindOnePromotion(promotionId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOnePromotionErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, house).Res()
}

func (h *promotionsHandlers) FindPromotion(c *fiber.Ctx) error {
	req := &promotions.PromotionFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findPromotionErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 100000
	}

	if req.OrderBy == "" {
		req.OrderBy = "created_at"
	}
	if req.Sort == "" {
		req.Sort = "DESC"
	}

	promotions := h.promotionsUsecase.FindPromotion(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, promotions).Res()
}

func (h *promotionsHandlers) AddPromotion(c *fiber.Ctx) error {
	req := &promotions.Promotion{
		HouseModel: make([]*promotions.PromotionHouseModel, 0),
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertPromotionErr),
			err.Error(),
		).Res()
	}

	if len(req.HouseModel) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertPromotionErr),
			"house model id is invalid",
		).Res()
	}

	promotion, err := h.promotionsUsecase.AddPromotion(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertPromotionErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, promotion).Res()
}

func (h *promotionsHandlers) UpdatePromotion(c *fiber.Ctx) error {
	promotionIdStr := strings.Trim(c.Params("promotion_id"), " ")
    promotionId, err := strconv.Atoi(promotionIdStr)
	if err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(updatePromotionErr),
            err.Error(),
        ).Res()
    }
	
	req := &promotions.Promotion{
		// HouseModel: make([]*promotions.PromotionHouseModel, 0),
		FreeItem: make([]*promotions.PromotionFreeItem, 0),
		Images:   make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updatePromotionErr),
			err.Error(),
		).Res()
	}
	req.Id = promotionId

	promotion, err := h.promotionsUsecase.UpdatePromotion(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updatePromotionErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, promotion).Res()
}

func (h *promotionsHandlers) DeletePromotion(c *fiber.Ctx) error {
	promotionId := strings.Trim(c.Params("promotion_id"), " ")

	promotion, err := h.promotionsUsecase.FindOnePromotion(promotionId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deletePromotionErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range promotion.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("promotions/%s", path.Base(p.Url)),
		})
	}

	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deletePromotionErr),
			err.Error(),
		).Res()
	}

	if err := h.promotionsUsecase.DeletePromotion(promotionId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deletePromotionErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
