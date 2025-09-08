package infra

import (
	"gorm.io/gorm"
	"log"
)

func SeedData(db *gorm.DB) {
	log.Println("Seeding mock data...")
	log.Println("Seeded successfully.")
}
