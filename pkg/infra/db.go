package infra

import (
	"database/sql"
	"fmt"
	"homemie/config"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func InitDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Port,
	)

	var db *sql.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				break // Success
			}
		}
		log.Printf("Failed to connect DB (attempt %d/5): %v", i+1, err)
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after multiple retries: %v", err)
	}

	log.Println("Database connected successfully")

	return db
}
