package analyticsRepositories

import (
	"fmt"
	"strconv"

	"github.com/yporn/sirarom-backend/modules/analytics"
	"google.golang.org/api/analyticsreporting/v4"
)

type IAnalyticsRepository interface {
	GetAnalyticsData() (*analytics.AnalyticsData, error)
}

type analyticsRepository struct {
	service *analyticsreporting.Service
	viewID  string
}

func AnalyticsRepository(service *analyticsreporting.Service, viewID string) (IAnalyticsRepository) {
	return &analyticsRepository{
		service: service, 
		viewID: viewID,
	}
}

// GetAnalyticsData fetches analytics data from Google Analytics
func (r *analyticsRepository) GetAnalyticsData() (*analytics.AnalyticsData, error) {
	// Construct the request
	request := &analyticsreporting.GetReportsRequest{
		ReportRequests: []*analyticsreporting.ReportRequest{
			{
				ViewId: r.viewID,
				DateRanges: []*analyticsreporting.DateRange{
					{
						StartDate: "7DaysAgo",
						EndDate:   "today",
					},
				},
				Metrics: []*analyticsreporting.Metric{
					{
						Expression: "ga:pageviews",
					},
					{
						Expression: "ga:uniquePageviews",
					},
				},
			},
		},
	}

	// Execute the request
	response, err := r.service.Reports.BatchGet(request).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analytics data: %v", err)
	}

	// Extract data from the response
	if len(response.Reports) == 0 || len(response.Reports[0].Data.Rows) == 0 {
		return nil, fmt.Errorf("no data returned from Google Analytics")
	}

	// Parse the data
	pageViews, err := strconv.Atoi(response.Reports[0].Data.Rows[0].Metrics[0].Values[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse page views: %v", err)
	}
	uniquePageViews, err := strconv.Atoi(response.Reports[0].Data.Rows[0].Metrics[1].Values[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse unique page views: %v", err)
	}

	// Create AnalyticsData object
	analyticsData := &analytics.AnalyticsData{
		PageViews:   pageViews,
		UniqueViews: uniquePageViews,
	}

	return analyticsData, nil
}
