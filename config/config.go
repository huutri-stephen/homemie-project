package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

)

type Config struct {
	Server struct {
		Port       string `yaml:"port"`
		Host       string `yaml:"host"`
		ApiVersion string `yaml:"version"`
	} `yaml:"app"`
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `env:"HOMIE_APP_DB_USER"`
		Password string `env:"HOMIE_APP_DB_PASSWORD"`
		Name     string `env:"HOMIE_APP_DB_NAME"`
		URL      string `yaml:"url"`
	} `yaml:"db"`
	Email struct {
		SmtpHost    string `env:"HOMIE_APP_SMTP_HOST"`
		SmtpPort    string `env:"HOMIE_APP_SMTP_PORT"`
		SmtpUser    string `env:"HOMIE_APP_SMTP_USER"`
		SmtpPass    string `env:"HOMIE_APP_SMTP_PASSWORD"`
		SenderEmail string `env:"HOMIE_APP_SMTP_SENDER"`
	}
	JWT struct {
		Secret string `env:"HOMIE_APP_JWT_SECRET"`
	}
	S3 struct {
		Endpoint         string `yaml:"endpoint"`
		ExternalEndpoint string `yaml:"external_endpoint"`
		AccessKey        string `env:"HOMIE_APP_S3_ACCESS_KEY_ID"`
		SecretKey        string `env:"HOMIE_APP_S3_SECRET_ACCESS_KEY"`
		Region           string `yaml:"region"`
		BucketName       string `yaml:"bucket_name"`
	} `yaml:"s3"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	if path == "" {
		path = "./config/config.yml"
	}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

