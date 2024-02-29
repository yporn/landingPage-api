package main

import (
	"os"

	"github.com/yporn/sirarom-backend/config"
	"github.com/yporn/sirarom-backend/modules/servers"
	"github.com/yporn/sirarom-backend/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())
	// fmt.Println(cfg.App())
	// fmt.Println(cfg.Jwt())
	db := databases.DbConnect(cfg.Db())
	defer db.Close()
	
	servers.NewServer(cfg, db).Start()
}