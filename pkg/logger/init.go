package logger

import (
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

func init() {
	var config zap.Config
	onGcpEnv := os.Getenv("ON_GCP")
	logLevel := os.Getenv("LOG_LEVEL")

	// if on the cloud we will use a production config
	if onGcpEnv == "true" {
		// create our uber zap configuration
		config = zapdriver.NewProductionConfig()
		// set the min logging level
		if logLevel == "" {
			config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		} else {
			var level zapcore.Level
			if err := level.UnmarshalText([]byte(logLevel)); err != nil {
				log.Printf("Invalid LOG_LEVEL '%s', defaulting to InfoLevel", logLevel)
				level = zap.InfoLevel
			}
			config.Level = zap.NewAtomicLevelAt(level)
		}
	} else {
		// running locally we will use a human-readable output
		config = zapdriver.NewDevelopmentConfig()
		config.Encoding = "console"
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// set the min logging level
		if logLevel == "" {
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		} else {
			var level zapcore.Level
			if err := level.UnmarshalText([]byte(logLevel)); err != nil {
				log.Printf("Invalid LOG_LEVEL '%s', defaulting to DebugLevel", logLevel)
				level = zap.DebugLevel
			}
			config.Level = zap.NewAtomicLevelAt(level)
		}
	}

	// creates our loggerInstance instance
	var err error = nil
	loggerInstance, err = config.Build()
	if err != nil {
		log.Fatalf("zap.config.Build(): %v", err)
	}
}
