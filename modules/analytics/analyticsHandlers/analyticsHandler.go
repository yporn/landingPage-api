package analyticsHandlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/analytics/analyticsUsecases"
	"github.com/yporn/sirarom-backend/modules/entities"
)

// IAnalyticsHandler defines the interface for the analytics handler
type IAnalyticsHandler interface {
	GetAnalyticsData(c *fiber.Ctx) error
}

// analyticsHandler implements the IAnalyticsHandler interface
type analyticsHandler struct {
	cfg            config.IConfig
	analyticsUsecase analyticsUsecases.IAnalyticsUsecase
}

// NewAnalyticsHandler creates a new instance of the analyticsHandler
func AnalyticsHandler(cfg config.IConfig, uc analyticsUsecases.IAnalyticsUsecase) IAnalyticsHandler {
	return &analyticsHandler{
		cfg:            cfg,
		analyticsUsecase: uc,
	}
}

// GetAnalyticsData retrieves analytics data
func (h *analyticsHandler) GetAnalyticsData(c *fiber.Ctx) error {
	// Call the use case to get analytics data
	data, err := h.analyticsUsecase.GetAnalyticsData()
	fmt.Println(data)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			"Failed to retrieve analytics data",
			err.Error(),
		).Res()
	}

	// Return the analytics data as JSON response
	return entities.NewResponse(c).Success(fiber.StatusOK, data).Res()
}
