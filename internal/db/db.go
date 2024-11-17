package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection and applies necessary configurations.
func InitDB() (*gorm.DB, error) {

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	// Retrieve database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" || dbSSLMode == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	// Form the connection string
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode,
	)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Retrieve the SQL database object
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting the SQL database object: %w", err)
	}

	// Ensure the UUID extension exists
	if _, err := sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"); err != nil {
		return nil, fmt.Errorf("error creating UUID extension: %w", err)
	}

	// Apply cascade delete constraints
	if err := applyCascadeDeleteConstraint(db); err != nil {
		return nil, fmt.Errorf("error applying cascade delete constraint: %w", err)
	}

	return db, nil
}

// CloseDB gracefully closes the database connection.
func CloseDB(dbInstance *gorm.DB) error {
	sqlDB, err := dbInstance.DB()
	if err != nil {
		return fmt.Errorf("error getting the SQL database object: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("error closing the database connection: %w", err)
	}
	return nil
}

// applyCascadeDeleteConstraint sets up the cascade delete constraints for the database.
func applyCascadeDeleteConstraint(db *gorm.DB) error {
	// Drop the existing foreign key constraint if it exists
	if err := db.Exec(`ALTER TABLE contact_infos DROP CONSTRAINT IF EXISTS fk_hotels_contact_infos;`).Error; err != nil {
		return fmt.Errorf("error dropping existing constraint: %w", err)
	}

	// Add the new foreign key constraint with cascade delete
	if err := db.Exec(`
        ALTER TABLE contact_infos
        ADD CONSTRAINT fk_hotels_contact_infos
        FOREIGN KEY (hotel_id)
        REFERENCES hotels(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE;
    `).Error; err != nil {
		return fmt.Errorf("error adding cascade delete constraint: %w", err)
	}

	return nil
}
