package hotel

import (
	"fmt"
	"hotel-guide/models"

	"github.com/google/uuid"
)

// HotelService provides operations on hotels
type HotelService struct {
	hotelRepo HotelRepository
}

// NewService creates a new instance of HotelService
func NewService(repo HotelRepository) *HotelService {
	return &HotelService{
		hotelRepo: repo,
	}
}

// CreateHotel creates a new hotel record
func (s *HotelService) CreateHotel(ownerName, ownerSurname, companyTitle string, contacts []models.ContactInfo) (*models.Hotel, error) {
	hotel := models.NewHotel(ownerName, ownerSurname, companyTitle, contacts)
	err := s.hotelRepo.Save(hotel) // Match the method name from repository (Save instead of CreateHotel)
	if err != nil {
		return nil, err
	}
	return hotel, nil
}

// DeleteHotel deletes an existing hotel by ID
func (s *HotelService) DeleteHotel(id uuid.UUID) error {
	err := s.hotelRepo.Delete(id) // Match method signature from repository (Delete instead of DeleteHotel)
	if err != nil {
		return fmt.Errorf("failed to delete hotel: %w", err)
	}
	return nil
}

// AddContactInfo adds a new contact info for the specified hotel
func (s *HotelService) AddContactInfo(hotelID uuid.UUID, contact models.ContactInfo) error {
	err := s.hotelRepo.AddContactInfo(hotelID, &contact) // Correct method name from repository
	if err != nil {
		return fmt.Errorf("failed to add contact info: %w", err)
	}
	return nil
}

// RemoveContactInfo removes a contact info from the specified hotel
func (s *HotelService) RemoveContactInfo(hotelID uuid.UUID, contactUUID uuid.UUID) error {
	err := s.hotelRepo.RemoveContactInfo(hotelID, contactUUID) // Correct method name from repository
	if err != nil {
		return fmt.Errorf("failed to remove contact info: %w", err)
	}
	return nil
}

// ListHotels lists all hotels
func (s *HotelService) ListHotels() ([]models.Hotel, error) {
	return s.hotelRepo.ListHotels() // Ensure ListHotels exists in your repository
}
