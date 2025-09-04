package infra

import (
	"fmt"
	"homemie/config"
	"homemie/models/dto"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.Config) *gorm.DB {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        cfg.DB.Host,
        cfg.DB.User,
        cfg.DB.Password,
        cfg.DB.Name,
        cfg.DB.Port,
    )

    var db *gorm.DB
    var err error

    for i := 0; i < 5; i++ { // thử 5 lần
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            break
        }
        log.Printf("Failed to connect DB (attempt %d/5): %v", i+1, err)
        time.Sleep(10 * time.Second)
    }
    
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    log.Println("Database connected successfully")

	err = db.AutoMigrate(
        &dto.EmailTemplate{},
    )

    if err != nil {
        log.Fatalf("AutoMigrate failed: %v", err)
    }

    log.Println("Database migrated successfully")
    return db
}
