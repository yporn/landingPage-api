package analyticsRepositories

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/yporn/sirarom-backend/modules/analytics"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

type IAnalyticsRepository interface {
	GetAnalyticsData() (*analytics.AnalyticsData, error)
}

type analyticsRepository struct {
	service    *analyticsreporting.Service
	propertyID string
}

func AnalyticsRepository(service *analyticsreporting.Service, propertyID string) IAnalyticsRepository {
	return &analyticsRepository{
		service:    service,
		propertyID: propertyID,
	}
}

// GetAnalyticsData fetches analytics data from Google Analytics.
func (r *analyticsRepository) GetAnalyticsData() (*analytics.AnalyticsData, error) {
	// Replace with your downloaded JSON key file path
	credPath := "credentials.json"

	// Read the JSON credentials file
	credData, err := ioutil.ReadFile(credPath)
	if err != nil {
		log.Fatalf("Error reading credentials file: %v", err)
	}

	// Create an authorized Analytics Data API service
	ctx := context.Background()
	svc, err := analyticsdata.NewService(ctx, option.WithCredentialsJSON(credData))
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}
	// Create the DimensionFilter
	filter := &analyticsdata.FilterExpression{
		Filter: &analyticsdata.Filter{
			FieldName: "pageTitle",
			StringFilter: &analyticsdata.StringFilter{
				MatchType: "EXACT",
				Value:     "/projects/1",
			},
		},
	}

	propertyID := "436823770"
	// Create the request
	req := &analyticsdata.RunReportRequest{
		Property: fmt.Sprintf("properties/%s", propertyID),
		DateRanges: []*analyticsdata.DateRange{
			{StartDate: "2024-03-28", EndDate: "today"},
		},
		Dimensions: []*analyticsdata.Dimension{
			{
				Name: "pageTitle",
			},
		},
		Metrics: []*analyticsdata.Metric{
			{
				Name: "screenPageViews",
			},
		},
		DimensionFilter: filter,
	}

	// Execute the request
	response, err := svc.Properties.RunReport("properties/436823770", req).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch analytics data: %v", err)
	}

	// Check if response.Rows is empty
	if len(response.Rows) == 0 {
		fmt.Println("No rows returned in the response.")
		return nil, nil // or return an appropriate error
	}

	// Access the first row and check if MetricValues is empty
	if len(response.Rows[0].MetricValues) == 0 {
		fmt.Println("No metric values found in the first row.")
		return nil, nil // or return an appropriate error
	}

	var pageViews int
	if len(response.Rows) > 0 && len(response.Rows[0].MetricValues) > 0 {
		pageViews, err = strconv.Atoi(response.Rows[0].MetricValues[0].Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse page views: %v", err)
		}
	} else {
		return nil, fmt.Errorf("no data available to parse")
	}

	fmt.Printf("Number of sessions today: %v\n", pageViews)

	// Create AnalyticsData object
	analyticsData := &analytics.AnalyticsData{
		PageViews: pageViews,
	}

	// Return the AnalyticsData object in JSON format
	return analyticsData, nil
}
