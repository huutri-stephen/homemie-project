package infra

import (
    "fmt"
    "log"
    "homemie/config"
	"homemie/models"

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

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    log.Println("Database connected successfully")

	err = db.AutoMigrate(
        &models.User{},
        &models.Listing{},
        &models.ListingImage{},
        &models.Favorite{},
        &models.Booking{},
        &models.Message{},
    )
    if err != nil {
        log.Fatalf("AutoMigrate failed: %v", err)
    }

    log.Println("Database migrated successfully")
    return db
}
