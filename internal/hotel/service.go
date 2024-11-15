package hotel

import (
	"fmt"

	"github.com/google/uuid"
)

type HotelService interface {
	CreateHotel(ownerName, ownerSurname, companyTitle string, contacts []ContactInfo) (*Hotel, error)
	DeleteHotel(id uuid.UUID) error
	AddContactInfo(hotelID uuid.UUID, contact *ContactInfo) error
	RemoveContactInfo(hotelID uuid.UUID, contactUUID uuid.UUID) error
	ListHotels() ([]Hotel, error)
	ListHotelOfficials() ([]HotelOfficial, error)
	GetHotelDetails(hotelID uuid.UUID) (*Hotel, error)
	FetchLocationStats(location string) (int, int, error)
}

// hotelService struct implements the HotelService interface
type hotelService struct {
	hotelRepo HotelRepository
}

func NewService(repo HotelRepository) HotelService {
	return &hotelService{
		hotelRepo: repo,
	}
}

func (s *hotelService) CreateHotel(ownerName, ownerSurname, companyTitle string, contacts []ContactInfo) (*Hotel, error) {
	hotel := NewHotel(ownerName, ownerSurname, companyTitle, contacts)
	if err := s.hotelRepo.Save(hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}

func (s *hotelService) DeleteHotel(id uuid.UUID) error {
	if err := s.hotelRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete hotel: %w", err)
	}
	return nil
}

func (s *hotelService) AddContactInfo(hotelID uuid.UUID, contact *ContactInfo) error {
	if err := s.hotelRepo.AddContactInfo(hotelID, contact); err != nil {
		return fmt.Errorf("failed to add contact info: %w", err)
	}
	return nil
}

func (s *hotelService) RemoveContactInfo(hotelID uuid.UUID, contactUUID uuid.UUID) error {
	if err := s.hotelRepo.RemoveContactInfo(hotelID, contactUUID); err != nil {
		return fmt.Errorf("failed to remove contact info: %w", err)
	}
	return nil
}

func (s *hotelService) ListHotels() ([]Hotel, error) {
	return s.hotelRepo.ListHotels()
}

func (s *hotelService) ListHotelOfficials() ([]HotelOfficial, error) {
	officials, err := s.hotelRepo.GetHotelOfficials()
	if err != nil {
		return nil, fmt.Errorf("failed to list hotel officials: %w", err)
	}
	return officials, nil
}

func (s *hotelService) GetHotelDetails(hotelID uuid.UUID) (*Hotel, error) {
	hotelDetails, err := s.hotelRepo.GetHotelDetails(hotelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel details: %w", err)
	}
	return hotelDetails, nil
}

func (s *hotelService) FetchLocationStats(location string) (int, int, error) {
	hotels, err := s.hotelRepo.FetchHotelsByLocation(location)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch hotels for location %s: %w", location, err)
	}

	hotelCount := len(hotels)
	phoneCount := 0
	for _, hotel := range hotels {
		for _, contact := range hotel.ContactInfos {
			if contact.InfoType == "phone" {
				phoneCount++
			}
		}
	}

	return hotelCount, phoneCount, nil
}
