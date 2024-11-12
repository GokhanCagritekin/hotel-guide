package hotel

import (
	"fmt"
	"hotel-guide/internal/db"
	"hotel-guide/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HotelRepository interface {
	Save(hotel *models.Hotel) error
	Delete(uuid uuid.UUID) error
	AddContactInfo(hotelUUID uuid.UUID, contact *models.ContactInfo) error
	RemoveContactInfo(hotelUUID, contactUUID uuid.UUID) error
	ListHotels() ([]models.Hotel, error)
	GetHotelOfficials() ([]models.HotelOfficial, error)
	GetHotelDetails(hotelID uuid.UUID) (*models.Hotel, error)
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
	return r.db.Where("id = ?", uuid).Delete(&models.Hotel{}).Error
}

func (r *hotelRepository) AddContactInfo(hotelUUID uuid.UUID, contact *models.ContactInfo) error {
	if contact.ID == uuid.Nil {
		contact.ID = uuid.New()
	}

	contact.HotelID = hotelUUID
	return r.db.Create(contact).Error
}

func (r *hotelRepository) RemoveContactInfo(hotelUUID, contactUUID uuid.UUID) error {
	var contact models.ContactInfo
	if err := r.db.Where("id = ? AND hotel_id = ?", contactUUID, hotelUUID).First(&contact).Error; err != nil {
		return fmt.Errorf("failed to find contact with ID %v for hotel with ID %v: %v", contactUUID, hotelUUID, err)
	}

	return r.db.Delete(&contact).Error
}

func (r *hotelRepository) ListHotels() ([]models.Hotel, error) {
	var hotels []models.Hotel
	err := r.db.Preload("ContactInfos").Find(&hotels).Error
	return hotels, err
}

func (r *hotelRepository) GetHotelOfficials() ([]models.HotelOfficial, error) {
	var officials []models.HotelOfficial
	err := r.db.Model(&models.Hotel{}).Select("owner_name, owner_surname, company_title").Find(&officials).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching hotel officials: %w", err)
	}
	return officials, nil
}

func (r *hotelRepository) GetHotelDetails(hotelID uuid.UUID) (*models.Hotel, error) {
	var hotel models.Hotel
	err := r.db.Preload("ContactInfos").First(&hotel, "id = ?", hotelID).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching hotel details: %w", err)
	}
	return &hotel, nil
}
