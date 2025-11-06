package main

import (
	"context"
	"flag"
	"homemie/config"
	"homemie/internal"
	"homemie/pkg/infra"
	"homemie/pkg/logger"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Info(context.Background(), "===== Start Homie Service ======")
	configPath := flag.String("CONFIG", "./config/config.yml", "Path of the config file")
	flag.Parse()

	logger.Infof(ctx, "Loading configuration from %s", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Errorf(ctx, "Failed to load configuration: %v", err)
	}

	db := infra.InitDB(cfg)

	// infra.SeedData(db)
	// infra.StartCronJobs(db, appLogger)

	// s3Client := infra.NewS3Client()
	// bucketName := cfg.S3.BucketName
	// infra.CreateBucketIfNotExists(s3Client, bucketName)
	
	r := internal.NewRouter(db, cfg, cfg.S3.ExternalEndpoint)

	// r := internal.NewRouter(db, cfg, appLogger, s3Client, cfg.S3.ExternalEndpoint)
	r.Run(":" + cfg.Server.Port)
}
