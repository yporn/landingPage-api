package generalHandlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/general"
	"github.com/yporn/sirarom-backend/modules/general/generalUsecases"
)

type generalHandlersErrCode string

const (
	findOneJGeneralErr generalHandlersErrCode = "general-001"
	updateGeneralErr   generalHandlersErrCode = "general-002"
)

type IGeneralHandler interface {
	FindOneGeneral(c *fiber.Ctx) error
	UpdateGeneral(c *fiber.Ctx) error
}

type generalHandler struct {
	cfg            config.IConfig
	generalUsecase generalUsecases.IGeneralUsecase
	filesUsecase   filesUsecases.IFilesUsecase
}

func GeneralHandler(cfg config.IConfig, generalUsecase generalUsecases.IGeneralUsecase, filesUsecase filesUsecases.IFilesUsecase) IGeneralHandler {
	return &generalHandler{
		cfg:            cfg,
		generalUsecase: generalUsecase,
		filesUsecase:   filesUsecase,
	}
}

func (h *generalHandler) FindOneGeneral(c *fiber.Ctx) error {
	generalId := strings.Trim(c.Params("general_id"), " ")

	general, err := h.generalUsecase.FindOneGeneral(generalId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneJGeneralErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, general).Res()
}

func (h *generalHandler) UpdateGeneral(c *fiber.Ctx) error {
	generalIdStr := strings.Trim(c.Params("general_id"), " ")
	generalId, err := strconv.Atoi(generalIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateGeneralErr),
			"Invalid job ID",
		).Res()
	}
	
	req := &general.General{
		Images: make([]*entities.Image, 0),
	}
	
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateGeneralErr),
			err.Error(),
		).Res()
	}

	req.Id = generalId

	job, err := h.generalUsecase.UpdateGeneral(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateGeneralErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, job).Res()
}
