package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Server struct {
        Port string
    }
    DB struct {
        Host     string
        Port     string
        User     string
        Password string
        Name     string
    }
}

func LoadConfig() Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system env")
    }

    cfg := Config{}
    cfg.Server.Port = getEnv("PORT", "8080")

    cfg.DB.Host = getEnv("DB_HOST", "localhost")
    cfg.DB.Port = getEnv("DB_PORT", "5432")
    cfg.DB.User = getEnv("DB_USER", "postgres")
    cfg.DB.Password = getEnv("DB_PASSWORD", "password")
    cfg.DB.Name = getEnv("DB_NAME", "mihome")

    return cfg
}

func getEnv(key, defaultVal string) string {
    val := os.Getenv(key)
    if val == "" {
        return defaultVal
    }
    return val
}
