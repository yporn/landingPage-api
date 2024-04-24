package seoHandlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/seo"
	"github.com/yporn/sirarom-backend/modules/seo/seoUsecases"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type seoHandlersErrCode string

const (
	findOneSeoErr seoHandlersErrCode = "seo-001"
	updateSeoErr  seoHandlersErrCode = "seo-002"
)

type ISeoHandler interface {
	FindOneSeo(c *fiber.Ctx) error
	UpdateSeo(c *fiber.Ctx) error
}

type seoHandler struct {
	cfg          config.IConfig
	seoUsecase   seoUsecases.ISeoUsecase
	db           *sql.DB
}

func SeoHandler(cfg config.IConfig, seoUsecase seoUsecases.ISeoUsecase, db *sql.DB) ISeoHandler {
	return &seoHandler{
		cfg:            cfg,
		seoUsecase: seoUsecase,
		db:             db,
	}
}

func (h *seoHandler) FindOneSeo(c *fiber.Ctx) error {
	seoId := strings.Trim(c.Params("seo_id"), " ")

	seo, err := h.seoUsecase.FindOneSeo(seoId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneSeoErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, seo).Res()
}

func (h *seoHandler) UpdateSeo(c *fiber.Ctx) error {
	seoIdStr := strings.Trim(c.Params("seo_id"), " ")
	seoId, err := strconv.Atoi(seoIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateSeoErr),
			"Invalid job ID",
		).Res()
	}

	req := &seo.Seo{}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateSeoErr),
			err.Error(),
		).Res()
	}

	req.Id = seoId

	job, err := h.seoUsecase.UpdateSeo(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateSeoErr),
			err.Error(),
		).Res()
	}

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "updated", "อัพเดตข้อมูล SEO")
	if err != nil {
		// Handle error if logging fails
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			fmt.Sprintf("Failed to log seo %v", userID),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, job).Res()
}

