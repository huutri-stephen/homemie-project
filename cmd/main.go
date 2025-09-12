package main

import (
	"homemie/config"
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

	s3Client := infra.NewS3Client()
	bucketName := cfg.S3.BucketName
	infra.CreateBucketIfNotExists(s3Client, bucketName)
	r := internal.NewRouter(db, cfg, appLogger, s3Client, cfg.S3.ExternalEndpoint)
	r.Run(":" + cfg.Server.Port)
}
