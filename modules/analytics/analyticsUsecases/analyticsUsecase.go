package analyticsUsecases

import (
	"github.com/yporn/sirarom-backend/modules/analytics"
	"github.com/yporn/sirarom-backend/modules/analytics/analyticsRepositories"
)

// IAnalyticsUsecase defines the interface for the analytics use case
type IAnalyticsUsecase interface {
	GetAnalyticsData() (*analytics.AnalyticsData, error)
}

// analyticsUsecase implements the IAnalyticsUsecase interface
type analyticsUsecase struct {
	analyticsRepository analyticsRepositories.IAnalyticsRepository
}

// NewAnalyticsUsecase creates a new instance of the analyticsUsecase
func AnalyticsUsecase(repo analyticsRepositories.IAnalyticsRepository) IAnalyticsUsecase {
	return &analyticsUsecase{analyticsRepository: repo}
}

// GetAnalyticsData retrieves analytics data
func (uc *analyticsUsecase) GetAnalyticsData() (*analytics.AnalyticsData, error) {
	return uc.analyticsRepository.GetAnalyticsData()
}
