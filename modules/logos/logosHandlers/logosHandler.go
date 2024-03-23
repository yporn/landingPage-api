package logosHandlers

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
	"github.com/yporn/sirarom-backend/modules/logos"
	"github.com/yporn/sirarom-backend/modules/logos/logosUsecases"
)

type logosHandlersErrCode string

const (
	findOneLogoErr logosHandlersErrCode = "logos-001"
	findLogoErr    logosHandlersErrCode = "logos-002"
	insertLogoErr  logosHandlersErrCode = "logos-003"
	deleteLogoErr  logosHandlersErrCode = "logos-004"
	updateLogoErr  logosHandlersErrCode = "logos-005"
)

type ILogosHandler interface {
	FindOneLogo(c *fiber.Ctx) error
	FindLogo(c *fiber.Ctx) error
	AddLogo(c *fiber.Ctx) error 
	UpdateLogo(c *fiber.Ctx) error 
	DeleteLogo(c *fiber.Ctx) error
}

type logosHandler struct {
	cfg             config.IConfig
	logosUsecase 	logosUsecases.ILogosUsecase
	filesUsecase    filesUsecases.IFilesUsecase
}

func LogosHandler(cfg config.IConfig, logosUsecase logosUsecases.ILogosUsecase, filesUsecase filesUsecases.IFilesUsecase) ILogosHandler {
	return &logosHandler{
		cfg:             cfg,
		logosUsecase: logosUsecase,
		filesUsecase:    filesUsecase,
	}
}

func (h *logosHandler) FindOneLogo(c *fiber.Ctx) error {
	logoId := strings.Trim(c.Params("logo_id"), " ")

	logo, err := h.logosUsecase.FindOneLogo(logoId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneLogoErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, logo).Res()
}

func (h *logosHandler) FindLogo(c *fiber.Ctx) error {
	req := &logos.LogoFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findLogoErr),
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

	logos := h.logosUsecase.FindLogo(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, logos).Res()
}

func (h *logosHandler) AddLogo(c *fiber.Ctx) error {
	req := &logos.Logo{
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertLogoErr),
			err.Error(),
		).Res()
	}

	logo, err := h.logosUsecase.AddLogo(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertLogoErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, logo).Res()
}

func (h *logosHandler) UpdateLogo(c *fiber.Ctx) error {
	logoIdStr := strings.Trim(c.Params("logo_id"), " ")
    logoId, err := strconv.Atoi(logoIdStr)
	if err != nil {
        return entities.NewResponse(c).Error(
            fiber.ErrBadRequest.Code,
            string(updateLogoErr),
            err.Error(),
        ).Res()
    }
	
	req := &logos.Logo{
		Images:   make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateLogoErr),
			err.Error(),
		).Res()
	}
	req.Id = logoId

	logo, err := h.logosUsecase.UpdateLogo(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateLogoErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, logo).Res()
}

func (h *logosHandler) DeleteLogo(c *fiber.Ctx) error {
	logoId := strings.Trim(c.Params("logo_id"), " ")

	logo, err := h.logosUsecase.FindOneLogo(logoId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteLogoErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range logo.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("logos/%s",  path.Base(p.Url)),
		})
	}

	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteLogoErr),
			err.Error(),
		).Res()
	}

	if err := h.logosUsecase.DeleteLogo(logoId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteLogoErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}