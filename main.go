package main

import (
	"context"
	"fmt"
	
	"log"
	"os"

	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/servers"
	"github.com/yporn/sirarom-backend/pkg/databases"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"

	"google.golang.org/api/option"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	// Replace with your Google Analytics 4 property ID
	propertyID := "436823770"

	// Replace this with the path to your service account credentials JSON file
	serviceAccountFile := "credentials.json"

	// Create a context.
	ctx := context.Background()

	// Initialize the Google Analytics Data client.
	client, err := analyticsdata.NewService(ctx, option.WithCredentialsFile(serviceAccountFile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	
	
	// Create the DimensionFilter
	filter := &analyticsdata.FilterExpression{
		Filter: &analyticsdata.Filter{
			FieldName: "pagePath",
			StringFilter: &analyticsdata.StringFilter{
				MatchType:     "BEGINS_WITH",
				Value:         "/projects/1/",
			},
		},
	}

	// Create the request
	request := &analyticsdata.RunReportRequest{
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
	response, err := client.Properties.RunReport("properties/436823770", request).Do()
	if err != nil {
		log.Fatalf("Failed to get report: %v", err)
	}

	// Print the report data
	for _, row := range response.Rows {
		fmt.Printf("Date: %s, Users: %v\n", row.DimensionValues[0].Value, row.MetricValues[0].Value)
	}

	cfg := config.LoadConfig(envPath())
	// fmt.Println(cfg.App())
	// fmt.Println(cfg.Jwt())
	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()

}
