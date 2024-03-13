package interestsHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/interests"
	"github.com/yporn/sirarom-backend/modules/interests/interestsUsecases"
)

type interestsHandlersErrCode string

const (
	findOneInterestErr interestsHandlersErrCode = "interests-001"
	// findProductErr    interestsHandlersErrCode = "products-002"
	insertInterestErr interestsHandlersErrCode = "interests-003"
	// deleteProductErr  interestsHandlersErrCode = "products-004"
	// updateProductErr  interestsHandlersErrCode = "products-005"
)

type IInterestsHandler interface {
	FindOneInterest(c *fiber.Ctx) error
	// FindProduct(c *fiber.Ctx) error
	AddInterest(c *fiber.Ctx) error
	// DeleteProduct(c *fiber.Ctx) error
	// UpdateProduct(c *fiber.Ctx) error
}

type interestsHandler struct {
	cfg              config.IConfig
	interestsUsecase interestsUsecases.IInterestsUsecase
	filesUsecase     filesUsecases.IFilesUsecase
}

func InterestsHandler(cfg config.IConfig, interestsUsecase interestsUsecases.IInterestsUsecase, filesUsecase filesUsecases.IFilesUsecase) IInterestsHandler {
	return &interestsHandler{
		cfg:              cfg,
		interestsUsecase: interestsUsecase,
		filesUsecase:     filesUsecase,
	}
}

func (h *interestsHandler) FindOneInterest(c *fiber.Ctx) error {
	interestId := strings.Trim(c.Params("interest_id"), " ")

	interest, err := h.interestsUsecase.FindOneInterest(interestId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneInterestErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, interest).Res()
}

func (h *interestsHandler) AddInterest(c *fiber.Ctx) error {
	req := &interests.Interest{}
		// Images: make([]*entities.Image, 0),
	
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertInterestErr),
			err.Error(),
		).Res()
	}

	interest, err := h.interestsUsecase.AddInterest(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertInterestErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, interest).Res()
}
