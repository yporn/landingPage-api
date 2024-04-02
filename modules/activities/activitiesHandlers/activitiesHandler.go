package activitiesHandlers

import (
	"database/sql"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/activities"
	"github.com/yporn/sirarom-backend/modules/activities/activitiesUsecases"
	"github.com/yporn/sirarom-backend/modules/entities"
	"github.com/yporn/sirarom-backend/modules/files"
	"github.com/yporn/sirarom-backend/modules/files/filesUsecases"
	"github.com/yporn/sirarom-backend/pkg/utils"
)

type activitiesHandlersErrCode string

const (
	findOneActivityErr activitiesHandlersErrCode = "activities-001"
	findActivityErr    activitiesHandlersErrCode = "activities-002"
	insertActivityErr  activitiesHandlersErrCode = "activities-003"
	deleteActivityErr  activitiesHandlersErrCode = "activities-004"
	updateActivityErr  activitiesHandlersErrCode = "activities-005"
)

type IActivitiesHandler interface {
	FindOneActivity(c *fiber.Ctx) error
	FindActivity(c *fiber.Ctx) error
	AddActivity(c *fiber.Ctx) error
	UpdateActivity(c *fiber.Ctx) error
	DeleteActivity(c *fiber.Ctx) error
}

type activitiesHandler struct {
	cfg               config.IConfig
	activitiesUsecase activitiesUsecases.IActivitiesUsecase
	filesUsecase      filesUsecases.IFilesUsecase
	db                *sql.DB
}

func ActivitiesHandler(cfg config.IConfig, activitiesUsecase activitiesUsecases.IActivitiesUsecase, filesUsecase filesUsecases.IFilesUsecase, db *sql.DB) IActivitiesHandler {
	return &activitiesHandler{
		cfg:               cfg,
		activitiesUsecase: activitiesUsecase,
		filesUsecase:      filesUsecase,
		db:                db,
	}
}

func (h *activitiesHandler) FindOneActivity(c *fiber.Ctx) error {
	activityId := strings.Trim(c.Params("activity_id"), " ")

	activity, err := h.activitiesUsecase.FindOneActivity(activityId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneActivityErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, activity).Res()
}

func (h *activitiesHandler) FindActivity(c *fiber.Ctx) error {
	req := &activities.ActivityFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findActivityErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 1000000
	}

	if req.OrderBy == "" {
		req.OrderBy = "created_at"
	}
	if req.Sort == "" {
		req.Sort = "DESC"
	}

	activities := h.activitiesUsecase.FindActivity(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, activities).Res()
}

func (h *activitiesHandler) AddActivity(c *fiber.Ctx) error {
	req := &activities.Activity{
		Images: make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertActivityErr),
			err.Error(),
		).Res()
	}

	activity, err := h.activitiesUsecase.AddActivity(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertActivityErr),
			err.Error(),
		).Res()
	}

	userID := utils.GetUserIDFromContext(c)
	// Log activity
	err = utils.LogActivity(h.db, strconv.Itoa(userID), "AddActivity", "Activity added: "+strconv.Itoa(activity.Id))
	if err != nil {
		// Handle error if logging fails
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			fmt.Sprintf("Failed to log activity %v",userID),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, activity).Res()
}

func (h *activitiesHandler) UpdateActivity(c *fiber.Ctx) error {
	activityIdStr := strings.Trim(c.Params("activity_id"), " ")
	activityId, err := strconv.Atoi(activityIdStr)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateActivityErr),
			err.Error(),
		).Res()
	}

	req := &activities.Activity{
		Images: make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateActivityErr),
			err.Error(),
		).Res()
	}
	req.Id = activityId

	activity, err := h.activitiesUsecase.UpdateActivity(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateActivityErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, activity).Res()
}

func (h *activitiesHandler) DeleteActivity(c *fiber.Ctx) error {
	activityId := strings.Trim(c.Params("activity_id"), " ")

	activity, err := h.activitiesUsecase.FindOneActivity(activityId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteActivityErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range activity.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("activities/%s", path.Base(p.Url)),
		})
	}

	if err := h.filesUsecase.DeleteFileOnStorage(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteActivityErr),
			err.Error(),
		).Res()
	}

	if err := h.activitiesUsecase.DeleteActivity(activityId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteActivityErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
