package activityLogsHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/activityLogs/activityLogsUsecases"
	"github.com/yporn/sirarom-backend/modules/entities"
)

type activityLogsHandlersErrCode string

const (
	findAllActivityLogsErr activityLogsHandlersErrCode = "activityLogs-001"
)

type IActivityLogsHandler interface {
	FindAllActivityLogs(c *fiber.Ctx) error
}

type activityLogsHandler struct {
	cfg                 config.IConfig
	activityLogsUsecases activityLogsUsecases.IActivityLogsUsecase
}

func ActivityLogsHandler(cfg config.IConfig, activityLogsUsecases activityLogsUsecases.IActivityLogsUsecase) IActivityLogsHandler {
	return &activityLogsHandler{
		cfg:                 cfg,
		activityLogsUsecases: activityLogsUsecases,
	}
}

func (h *activityLogsHandler) FindAllActivityLogs(c *fiber.Ctx) error {
	activityLogs, err := h.activityLogsUsecases.FindAllActivityLogs()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findAllActivityLogsErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, activityLogs).Res()
}