package hotel

import (
	"hotel-guide/internal/db"
	"hotel-guide/models"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type HotelRepository interface {
	Save(hotel *models.Hotel) error
	Delete(uuid uuid.UUID) error
	AddContactInfo(hotelUUID uuid.UUID, contact *models.ContactInfo) error
	RemoveContactInfo(hotelUUID, contactUUID uuid.UUID) error
	ListHotels() ([]models.Hotel, error) // Add ListHotels method to your repository interface
}

type hotelRepository struct {
	db *gorm.DB
}

func NewRepository() HotelRepository {
	return &hotelRepository{db: db.DB}
}

func (r *hotelRepository) Save(hotel *models.Hotel) error {
	return r.db.Create(hotel).Error
}

func (r *hotelRepository) Delete(uuid uuid.UUID) error {
	return r.db.Where("uuid = ?", uuid).Delete(&models.Hotel{}).Error
}

func (r *hotelRepository) AddContactInfo(hotelUUID uuid.UUID, contact *models.ContactInfo) error {
	return r.db.Model(&models.Hotel{}).Where("uuid = ?", hotelUUID).Association("ContactInfos").Append(contact).Error
}

func (r *hotelRepository) RemoveContactInfo(hotelUUID, contactUUID uuid.UUID) error {
	var contact models.ContactInfo
	err := r.db.Where("uuid = ?", contactUUID).First(&contact).Error
	if err != nil {
		return err
	}
	return r.db.Model(&models.Hotel{}).Where("uuid = ?", hotelUUID).Association("ContactInfos").Delete(contact).Error
}

func (r *hotelRepository) ListHotels() ([]models.Hotel, error) {
	var hotels []models.Hotel
	err := r.db.Find(&hotels).Error
	return hotels, err
}
