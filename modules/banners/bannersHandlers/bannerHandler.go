package bannersHandlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/banners"
	"github.com/yporn/sirarom-backend/modules/banners/bannersUsecases"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type bannersHandlersErrCode string

const (
	findOneBannerErr bannersHandlersErrCode = "banners-001"
	findBannerErr    bannersHandlersErrCode = "banners-002"
	insertBannerErr  bannersHandlersErrCode = "banners-003"
	deleteBannerErr  bannersHandlersErrCode = "banners-004"
	updateBannerErr  bannersHandlersErrCode = "banners-005"
)

type IBannersHandler interface {
	FindOneBanner(c *fiber.Ctx) error
	FindBanner(c *fiber.Ctx) error
	AddBanner(c *fiber.Ctx) error
	UpdateBanner(c *fiber.Ctx) error
	DeleteBanner(c *fiber.Ctx) error
}

type bannersHandler struct {
	cfg            config.IConfig
	bannersUsecase bannersUsecases.IBannersUsecase
	filesUsecase   filesUsecases.IFilesUsecase
	db             *sql.DB
}

func BannersHandler(cfg config.IConfig, bannersUsecase bannersUsecases.IBannersUsecase, filesUsecase filesUsecases.IFilesUsecase, db *sql.DB) IBannersHandler {
	return &bannersHandler{
		cfg:            cfg,
		bannersUsecase: bannersUsecase,
		filesUsecase:   filesUsecase,
		db:             db,
	}
}

func (h *bannersHandler) FindOneBanner(c *fiber.Ctx) error {
	bannerId := strings.Trim(c.Params("banner_id"), " ")

	banner, err := h.bannersUsecase.FindOneBanner(bannerId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneBannerErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, banner).Res()
}

func (h *bannersHandler) FindBanner(c *fiber.Ctx) error {
	req := &banners.BannerFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findBannerErr),
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
		req.OrderBy = "index"
	}
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	banners := h.bannersUsecase.FindBanner(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, banners).Res()
}

func (h *bannersHandler) AddBanner(c *fiber.Ctx) error {
	req := &banners.Banner{
		Images: make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertBannerErr),
			err.Error(),
		).Res()
	}

	banner, err := h.bannersUsecase.AddBanner(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertBannerErr),
			err.Error(),
		).Res()
	}

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "created", "เพิ่มข้อมูลแบนเนอร์")
	if err != nil {
		// Handle error if logging fails
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			fmt.Sprintf("Failed to log activity %v", userID),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, banner).Res()
}

func (h *bannersHandler) UpdateBanner(c *fiber.Ctx) error {
	bannerIdStr := strings.Trim(c.Params("banner_id"), " ")
	bannerId, err := strconv.Atoi(bannerIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateBannerErr),
			err.Error(),
		).Res()
	}

	req := &banners.Banner{
		Images: make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateBannerErr),
			err.Error(),
		).Res()
	}
	req.Id = bannerId

	banner, err := h.bannersUsecase.UpdateBanner(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateBannerErr),
			err.Error(),
		).Res()
	}

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "updated", "แก้ไขข้อมูลแบนเนอร์")
	if err != nil {
		// Handle error if logging fails
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			fmt.Sprintf("Failed to log activity %v", userID),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, banner).Res()
}

func (h *bannersHandler) DeleteBanner(c *fiber.Ctx) error {
	bannerId := strings.Trim(c.Params("banner_id"), " ")

	banner, err := h.bannersUsecase.FindOneBanner(bannerId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteBannerErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range banner.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("banner/%s", p.FileName),
		})
	}

	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteBannerErr),
			err.Error(),
		).Res()
	}

	if err := h.bannersUsecase.DeleteBanner(bannerId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteBannerErr),
			err.Error(),
		).Res()
	}

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "deleted", "ลบข้อมูลแบนเนอร์")
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
