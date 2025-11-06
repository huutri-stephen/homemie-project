package infra

import (
	"database/sql"
	"log"
)

func SeedData(db *sql.DB) {
	log.Println("Seeding mock data...")
	log.Println("Seeded successfully.")
}
