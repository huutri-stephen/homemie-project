package main

import (
	"homemie/config"
	"homemie/internal"
	"homemie/internal/infra"
)

func main() {
	cfg := config.LoadConfig()
	db := infra.InitDB(cfg)

	infra.SeedData(db)
	infra.CronJobs(db)

	r := internal.NewRouter(db, cfg)
	r.Run(":" + cfg.Server.Port)
}