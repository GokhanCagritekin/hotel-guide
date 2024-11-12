package hotel

import (
	"fmt"
	"hotel-guide/models"

	"github.com/google/uuid"
)

type HotelService struct {
	hotelRepo HotelRepository
}

func NewService(repo HotelRepository) *HotelService {
	return &HotelService{
		hotelRepo: repo,
	}
}

func (s *HotelService) CreateHotel(ownerName, ownerSurname, companyTitle string, contacts []models.ContactInfo) (*models.Hotel, error) {
	hotel := models.NewHotel(ownerName, ownerSurname, companyTitle, contacts)
	if err := s.hotelRepo.Save(hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}

func (s *HotelService) DeleteHotel(id uuid.UUID) error {
	if err := s.hotelRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete hotel: %w", err)
	}
	return nil
}

func (s *HotelService) AddContactInfo(hotelID uuid.UUID, contact *models.ContactInfo) error {
	if err := s.hotelRepo.AddContactInfo(hotelID, contact); err != nil {
		return fmt.Errorf("failed to add contact info: %w", err)
	}
	return nil
}

func (s *HotelService) RemoveContactInfo(hotelID uuid.UUID, contactUUID uuid.UUID) error {
	if err := s.hotelRepo.RemoveContactInfo(hotelID, contactUUID); err != nil {
		return fmt.Errorf("failed to remove contact info: %w", err)
	}
	return nil
}

func (s *HotelService) ListHotels() ([]models.Hotel, error) {
	return s.hotelRepo.ListHotels()
}

func (s *HotelService) ListHotelOfficials() ([]models.HotelOfficial, error) {
	officials, err := s.hotelRepo.GetHotelOfficials()
	if err != nil {
		return nil, fmt.Errorf("failed to list hotel officials: %w", err)
	}
	return officials, nil
}

func (s *HotelService) GetHotelDetails(hotelID uuid.UUID) (*models.Hotel, error) {
	hotelDetails, err := s.hotelRepo.GetHotelDetails(hotelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel details: %w", err)
	}
	return hotelDetails, nil
}
