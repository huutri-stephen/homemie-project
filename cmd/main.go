package main

import (
	"github.com/Stephen-R-L/homemie-project/config"
	"homemie/internal"
	"homemie/internal/infra"
	"homemie/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Logger
	appLogger := logger.InitLogger("homemie-project")
	defer appLogger.Sync() // Flushes buffer, if any

	db := infra.InitDB(cfg)

	infra.SeedData(db)
	infra.StartCronJobs(db, appLogger)

	r := internal.NewRouter(db, cfg, appLogger)
	r.Run(":" + cfg.Server.Port)
}
