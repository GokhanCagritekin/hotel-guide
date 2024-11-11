package db

import (
	"hotel-guide/models"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDB() {
	var err error
	connectionString := "host=localhost user=youruser dbname=hotel_guide password=yourpassword sslmode=disable"

	DB, err = gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Veritabanı otomatik migration'ı
	DB.AutoMigrate(&models.Hotel{}, &models.Report{}, &models.ContactInfo{})
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error closing the database: %v", err)
	}
}
