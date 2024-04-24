package main

import (
	// "context"
	// "fmt"
	// "log"
	// "context"
	// "fmt"
	// "io/ioutil"
	// "log"
	"os"

	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/servers"
	"github.com/yporn/sirarom-backend/pkg/databases"
	// analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	// "google.golang.org/api/option"
	// analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	// "google.golang.org/api/option"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	// // Replace with your downloaded JSON key file path
	// credPath := "credentials.json"

	// // Read the JSON credentials file
	// credData, err := ioutil.ReadFile(credPath)
	// if err != nil {
	// 	log.Fatalf("Error reading credentials file: %v", err)
	// }

	// // Create an authorized Analytics Data API service
	// ctx := context.Background()
	// svc, err := analyticsdata.NewService(ctx, option.WithCredentialsJSON(credData))
	// if err != nil {
	// 	log.Fatalf("Error creating service: %v", err)
	// }

	// // Define the request
	// req := &analyticsdata.RunReportRequest{
	// 	Property: "properties/436823770",
	// 	DateRanges: []*analyticsdata.DateRange{
	// 		{StartDate: "yesterday", EndDate: "today"},
	// 	},
	// 	Dimensions: []*analyticsdata.Dimension{
	// 		{Name: "date"},
	// 	},
	// 	Metrics: []*analyticsdata.Metric{
	// 		{Name: "activeUsers"},
	// 	},
	// }

	// // Run the report
	// resp, err := svc.Properties.RunReport("properties/436823770", req).Do()
	// if err != nil {
	// 	log.Fatalf("Error running report: %v", err)
	// }

	// // Print the report data
	// if len(resp.Rows) > 0 {
	// 	fmt.Printf("Number of sessions today: %v\n", resp.Rows[0].MetricValues[0].Value)
	// } else {
	// 	fmt.Println("No data found for today.")
	// }

    
	cfg := config.LoadConfig(envPath())
	// fmt.Println(cfg.App())
	// fmt.Println(cfg.Jwt())
	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()

}
