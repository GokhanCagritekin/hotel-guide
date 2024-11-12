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
	// Eğer contact ID boşsa, yeni bir UUID atayın
	if contact.ID == uuid.Nil {
		contact.ID = uuid.New()
	}

	// ContactInfo'yu doğrudan veritabanına ekleyin
	contact.HotelID = hotelUUID
	return r.db.Create(contact).Error
}

func (r *hotelRepository) RemoveContactInfo(hotelUUID, contactUUID uuid.UUID) error {
	var contact models.ContactInfo

	// Find the ContactInfo by UUID and ensure it belongs to the hotel
	if err := r.db.Where("id = ? AND hotel_id = ?", contactUUID, hotelUUID).First(&contact).Error; err != nil {
		// Return a more informative error if contact is not found or does not belong to the hotel
		return fmt.Errorf("failed to find contact with ID %v for hotel with ID %v: %v", contactUUID, hotelUUID, err)
	}

	// Delete the ContactInfo record
	if err := r.db.Delete(&contact).Error; err != nil {
		// Return an error with more context if deletion fails
		return fmt.Errorf("failed to delete contact with ID %v for hotel with ID %v: %v", contactUUID, hotelUUID, err)
	}

	return nil
}

func (r *hotelRepository) ListHotels() ([]models.Hotel, error) {
	var hotels []models.Hotel
	// Gerekli ilişkiyi yüklemek için Preload kullanıyoruz
	err := r.db.Preload("ContactInfos").Find(&hotels).Error
	return hotels, err
}

func (r *hotelRepository) GetHotelOfficials() ([]models.HotelOfficial, error) {
	var officials []models.HotelOfficial
	// Query the database for hotel officials (no hotelID filter)
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
