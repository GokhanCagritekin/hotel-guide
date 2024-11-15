package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Retrieve database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Form the connection string
	connectionString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + dbSSLMode

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting the database object: %w", err)
	}

	sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := applyCascadeDeleteConstraint(db); err != nil {
		return nil, fmt.Errorf("error applying cascade delete constraint: %w", err)
	}

	return db, nil
}

func CloseDB(dbInstance *gorm.DB) {
	sqlDB, err := dbInstance.DB()
	if err != nil {
		log.Fatalf("Error getting the database object: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Error closing the database: %v", err)
	}
}

func applyCascadeDeleteConstraint(db *gorm.DB) error {
	if err := db.Exec(`ALTER TABLE contact_infos DROP CONSTRAINT IF EXISTS fk_hotels_contact_infos`).Error; err != nil {
		return err
	}

	return db.Exec(`
        ALTER TABLE contact_infos
        ADD CONSTRAINT fk_hotels_contact_infos
        FOREIGN KEY (hotel_id)
        REFERENCES hotels(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE;
    `).Error
}
