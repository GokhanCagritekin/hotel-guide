package hotel

import (
	"fmt"

	"github.com/google/uuid"
)

// Service provides operations on hotels
type Service struct {
	hotelRepo HotelRepository
}

// HotelRepository defines the methods required to interact with the data storage
type HotelRepository interface {
	CreateHotel(hotel *Hotel) error
	DeleteHotel(id uuid.UUID) error
	AddContactInfo(hotelID uuid.UUID, contact ContactInfo) error
	RemoveContactInfo(hotelID uuid.UUID, contact ContactInfo) error
	ListHotels() ([]Hotel, error)
}

// NewService creates a new instance of HotelService
func NewService(repo HotelRepository) *Service {
	return &Service{
		hotelRepo: repo,
	}
}

// CreateHotel creates a new hotel record
func (s *Service) CreateHotel(ownerName, ownerSurname, companyTitle string, contacts []ContactInfo) (*Hotel, error) {
	hotel := NewHotel(ownerName, ownerSurname, companyTitle, contacts)
	err := s.hotelRepo.CreateHotel(hotel)
	if err != nil {
		return nil, err
	}
	return hotel, nil
}

// DeleteHotel deletes an existing hotel by ID
func (s *Service) DeleteHotel(id uuid.UUID) error {
	err := s.hotelRepo.DeleteHotel(id)
	if err != nil {
		return fmt.Errorf("failed to delete hotel: %w", err)
	}
	return nil
}

// AddContactInfo adds a new contact info for the specified hotel
func (s *Service) AddContactInfo(hotelID uuid.UUID, contact ContactInfo) error {
	err := s.hotelRepo.AddContactInfo(hotelID, contact)
	if err != nil {
		return fmt.Errorf("failed to add contact info: %w", err)
	}
	return nil
}

// RemoveContactInfo removes a contact info from the specified hotel
func (s *Service) RemoveContactInfo(hotelID uuid.UUID, contact ContactInfo) error {
	err := s.hotelRepo.RemoveContactInfo(hotelID, contact)
	if err != nil {
		return fmt.Errorf("failed to remove contact info: %w", err)
	}
	return nil
}

// ListHotels lists all hotels
func (s *Service) ListHotels() ([]Hotel, error) {
	return s.hotelRepo.ListHotels()
}
