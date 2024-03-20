package promotionsHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
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
}

type promotionsHandlers struct {
	cfg               config.IConfig
	promotionsUsecase promotionsUsecases.IPromotionsUsecase
	filesUsecase      filesUsecases.IFilesUsecase
}

func PromotionsHandlers(cfg config.IConfig, promotionsUsecase promotionsUsecases.IPromotionsUsecase, filesUsecase filesUsecases.IFilesUsecase) IPromotionsHandler {
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

