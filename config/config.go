package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port       string
		Host       string
		ApiVersion string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Email struct {
		SmtpHost    string
		SmtpPort    string
		SmtpUser    string
		SmtpPass    string
		SenderEmail string
	}
	S3 struct {
		Endpoint   string
		AccessKey  string
		SecretKey  string
		Region     string
		BucketName string
	}
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	cfg := Config{}
	cfg.Server.Host = getEnv("APP_HOST", "localhost")
	cfg.Server.Port = getEnv("APP_PORT", "8080")
	cfg.Server.ApiVersion = getEnv("API_VERSION", "/api/v1")

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres")
	cfg.DB.Password = getEnv("DB_PASSWORD", "password")
	cfg.DB.Name = getEnv("DB_NAME", "homemie")

	cfg.Email.SmtpHost = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.Email.SmtpPort = getEnv("SMTP_PORT", "587")
	cfg.Email.SmtpUser = getEnv("SMTP_USER", "noreply.homemie@gmail.com")
	cfg.Email.SmtpPass = getEnv("SMTP_PASSWORD", "your_smtp_password")
	cfg.Email.SenderEmail = getEnv("SMTP_SENDER", "noreply.homemie@gmail.com")

	cfg.S3.Endpoint = getEnv("S3_ENDPOINT", "http://localhost:9000")
	cfg.S3.AccessKey = getEnv("S3_ACCESS_KEY_ID", "admin")
	cfg.S3.SecretKey = getEnv("S3_SECRET_ACCESS_KEY", "admin123")
	cfg.S3.Region = getEnv("S3_REGION", "us-east-1")
	cfg.S3.BucketName = getEnv("S3_BUCKET_NAME", "homemie-media")

	return cfg
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
