package db

import (
	"hotel-guide/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	connectionString := "postgres://myuser:mysecretpassword@localhost:5432/hotels?sslmode=disable"

	// Open connection using the gorm.io/driver/postgres driver
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// UUID extension (use raw SQL to enable UUID extension if needed)
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting the database object: %v", err)
	}
	// Executing SQL to enable UUID extension
	sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Veritabanı otomatik migration'ı
	if err := DB.AutoMigrate(&models.Hotel{}, &models.Report{}, &models.ContactInfo{}); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	if err := applyCascadeDeleteConstraint(DB); err != nil {
		log.Fatalf("Error applying cascade delete constraint: %v", err)
	}
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting the database object: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Error closing the database: %v", err)
	}
}

func applyCascadeDeleteConstraint(db *gorm.DB) error {
	// Drop the existing foreign key constraint if it exists
	if err := db.Exec(`ALTER TABLE contact_infos DROP CONSTRAINT IF EXISTS fk_hotels_contact_infos`).Error; err != nil {
		return err
	}

	// Add the new foreign key constraint with ON DELETE CASCADE and ON UPDATE CASCADE
	return db.Exec(`
        ALTER TABLE contact_infos
        ADD CONSTRAINT fk_hotels_contact_infos
        FOREIGN KEY (hotel_id)
        REFERENCES hotels(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE;
    `).Error
}
