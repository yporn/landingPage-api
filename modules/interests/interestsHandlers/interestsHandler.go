package interestsHandlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/modules/interests"
	"github.com/yporn/sirarom-backend/modules/interests/interestsUsecases"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type interestsHandlersErrCode string

const (
	findOneInterestErr interestsHandlersErrCode = "interests-001"
	findInterestErr    interestsHandlersErrCode = "interests-002"
	insertInterestErr  interestsHandlersErrCode = "interests-003"
	deleteInterestErr  interestsHandlersErrCode = "interests-004"
	updateInterestErr  interestsHandlersErrCode = "interests-005"
)

type IInterestsHandler interface {
	FindOneInterest(c *fiber.Ctx) error
	FindInterest(c *fiber.Ctx) error
	AddInterest(c *fiber.Ctx) error
	DeleteInterest(c *fiber.Ctx) error
	UpdateInterest(c *fiber.Ctx) error
}

type interestsHandler struct {
	cfg              config.IConfig
	interestsUsecase interestsUsecases.IInterestsUsecase
	filesUsecase     filesUsecases.IFilesUsecase
	db               *sql.DB
}

func InterestsHandler(cfg config.IConfig, interestsUsecase interestsUsecases.IInterestsUsecase, filesUsecase filesUsecases.IFilesUsecase, db *sql.DB) IInterestsHandler {
	return &interestsHandler{
		cfg:              cfg,
		interestsUsecase: interestsUsecase,
		filesUsecase:     filesUsecase,
		db:               db,
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

func (h *interestsHandler) FindInterest(c *fiber.Ctx) error {
	req := &interests.InterestFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findInterestErr),
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
		req.OrderBy = "id"
	}
	if req.Sort == "" {
		req.Sort = "DESC"
	}

	interests := h.interestsUsecase.FindInterest(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, interests).Res()
}

func (h *interestsHandler) AddInterest(c *fiber.Ctx) error {
	req := &interests.Interest{
		Images: make([]*entities.Image, 0),
	}

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

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "created", "เพิ่มข้อมูลดอกเบี้ย : "+interest.BankName)
	if err != nil {
		// Handle error if logging fails
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			fmt.Sprintf("Failed to log activity %v", userID),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, interest).Res()
}

func (h *interestsHandler) DeleteInterest(c *fiber.Ctx) error {
	interestId := strings.Trim(c.Params("interest_id"), " ")

	// Check if the interest ID is empty
	if interestId == "" {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteInterestErr),
			"Interest ID is empty",
		).Res()
	}
	interest, err := h.interestsUsecase.FindOneInterest(interestId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteInterestErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range interest.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("interests/%s", p.FileName),
		})
	}

	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteInterestErr),
			err.Error(),
		).Res()
	}

	if err := h.interestsUsecase.DeleteInterest(interestId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteInterestErr),
			err.Error(),
		).Res()
	}

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "deleted", "ลบข้อมูลดอกเบี้ย : "+interest.BankName)
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

func (h *interestsHandler) UpdateInterest(c *fiber.Ctx) error {
	interestIdStr := strings.Trim(c.Params("interest_id"), " ")
	interestId, err := strconv.Atoi(interestIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateInterestErr),
			err.Error(),
		).Res()
	}

	req := &interests.Interest{
		Images: make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateInterestErr),
			err.Error(),
		).Res()
	}
	req.Id = interestId

	interest, err := h.interestsUsecase.UpdateInterest(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateInterestErr),
			err.Error(),
		).Res()
	}

	// Log activity
	userID := utils.GetUserIDFromContext(c)
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "updated", "แก้ไขข้อมูลดอกเบี้ย : "+interest.BankName)
	if err != nil {
		// Handle error if logging fails
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			fmt.Sprintf("Failed to log activity %v", userID),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, interest).Res()
}
