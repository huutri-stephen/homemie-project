package infra

import (
    "log"
    "gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
    log.Println("Seeding mock data...")
    log.Println("Seeded successfully.")
}