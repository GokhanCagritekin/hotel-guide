package hotel

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HotelRepository interface {
	Save(hotel *Hotel) error
	Delete(uuid uuid.UUID) error
	AddContactInfo(hotelUUID uuid.UUID, contact *ContactInfo) error
	RemoveContactInfo(hotelUUID, contactUUID uuid.UUID) error
	ListHotels() ([]Hotel, error)
	GetHotelOfficials() ([]HotelOfficial, error)
	GetHotelDetails(hotelID uuid.UUID) (*Hotel, error)
	FetchHotelsByLocation(location string) ([]Hotel, error)
}

type hotelRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{db: db}
}

func (r *hotelRepository) Save(hotel *Hotel) error {
	return r.db.Create(hotel).Error
}

func (r *hotelRepository) Delete(uuid uuid.UUID) error {
	return r.db.Where("id = ?", uuid).Delete(&Hotel{}).Error
}

func (r *hotelRepository) AddContactInfo(hotelUUID uuid.UUID, contact *ContactInfo) error {
	if contact.ID == uuid.Nil {
		contact.ID = uuid.New()
	}

	contact.HotelID = hotelUUID
	return r.db.Create(contact).Error
}

func (r *hotelRepository) RemoveContactInfo(hotelUUID, contactUUID uuid.UUID) error {
	var contact ContactInfo
	if err := r.db.Where("id = ? AND hotel_id = ?", contactUUID, hotelUUID).First(&contact).Error; err != nil {
		return fmt.Errorf("failed to find contact with ID %v for hotel with ID %v: %v", contactUUID, hotelUUID, err)
	}

	return r.db.Delete(&contact).Error
}

func (r *hotelRepository) ListHotels() ([]Hotel, error) {
	var hotels []Hotel
	err := r.db.Preload("ContactInfos").Find(&hotels).Error
	return hotels, err
}

func (r *hotelRepository) GetHotelOfficials() ([]HotelOfficial, error) {
	var officials []HotelOfficial
	err := r.db.Model(&Hotel{}).Select("owner_name, owner_surname, company_title").Find(&officials).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching hotel officials: %w", err)
	}
	return officials, nil
}

func (r *hotelRepository) GetHotelDetails(hotelID uuid.UUID) (*Hotel, error) {
	var hotel Hotel
	err := r.db.Preload("ContactInfos").First(&hotel, "id = ?", hotelID).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching hotel details: %w", err)
	}
	return &hotel, nil
}

func (r *hotelRepository) FetchHotelsByLocation(location string) ([]Hotel, error) {
	var hotels []Hotel

	err := r.db.Joins("JOIN contact_infos ON contact_infos.hotel_id = hotels.id").
		Where("contact_infos.info_type IN (?)", []string{"location", "phone"}).
		Where("contact_infos.info_content = ?", location).
		Preload("ContactInfos").
		Find(&hotels).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching hotels by location %s: %w", location, err)
	}
	return hotels, nil
}
