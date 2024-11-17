package hotel

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHotelRepository struct {
	mock.Mock
}

func (m *MockHotelRepository) Save(hotel *Hotel) error {
	args := m.Called(hotel)
	return args.Error(0)
}

func (m *MockHotelRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHotelRepository) AddContactInfo(hotelID uuid.UUID, contact *ContactInfo) error {
	args := m.Called(hotelID, contact)
	return args.Error(0)
}

func (m *MockHotelRepository) RemoveContactInfo(hotelID uuid.UUID, contactUUID uuid.UUID) error {
	args := m.Called(hotelID, contactUUID)
	return args.Error(0)
}

func (m *MockHotelRepository) ListHotels() ([]Hotel, error) {
	args := m.Called()
	return args.Get(0).([]Hotel), args.Error(1)
}

func (m *MockHotelRepository) GetHotelOfficials() ([]HotelOfficial, error) {
	args := m.Called()
	return args.Get(0).([]HotelOfficial), args.Error(1)
}

func (m *MockHotelRepository) GetHotelDetails(hotelID uuid.UUID) (*Hotel, error) {
	args := m.Called(hotelID)
	return args.Get(0).(*Hotel), args.Error(1)
}

func (m *MockHotelRepository) FetchHotelsByLocation(location string) ([]Hotel, error) {
	args := m.Called(location)
	return args.Get(0).([]Hotel), args.Error(1)
}

func TestCreateHotel(t *testing.T) {
	// Mock repository creation
	mockRepo := new(MockHotelRepository)

	// Define test hotel
	hotel := &Hotel{
		OwnerName:    "John",
		OwnerSurname: "Doe",
		CompanyTitle: "Doe Ltd.",
	}

	// Expect the Save method to be called once with the correct arguments
	mockRepo.On("Save", mock.MatchedBy(func(h *Hotel) bool {
		// Check if the ID is not nil and the name/surname match the test case
		return h.ID != uuid.Nil && h.OwnerName == "John" && h.OwnerSurname == "Doe"
	})).Return(nil).Once()

	// Create the service with the mocked repository
	service := NewService(mockRepo)

	// Call CreateHotel
	createdHotel, err := service.CreateHotel(hotel.OwnerName, hotel.OwnerSurname, hotel.CompanyTitle, nil)

	// Assert no error occurred and the hotel was created with the expected values
	assert.NoError(t, err)
	assert.Equal(t, hotel.OwnerName, createdHotel.OwnerName)
	assert.Equal(t, hotel.OwnerSurname, createdHotel.OwnerSurname)
	assert.Equal(t, hotel.CompanyTitle, createdHotel.CompanyTitle)

	// Ensure that the Save method was called as expected
	mockRepo.AssertExpectations(t)
}

func TestDeleteHotel(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	hotelID := uuid.New()

	mockRepo.On("Delete", hotelID).Return(nil).Once()

	err := service.DeleteHotel(hotelID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAddContactInfo(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	hotelID := uuid.New()
	contact := &ContactInfo{
		InfoType:    ContactTypePhone,
		InfoContent: "123-456-7890",
	}

	mockRepo.On("AddContactInfo", hotelID, contact).Return(nil).Once()

	err := service.AddContactInfo(hotelID, contact)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestRemoveContactInfo(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	hotelID := uuid.New()
	contactID := uuid.New()

	mockRepo.On("RemoveContactInfo", hotelID, contactID).Return(nil).Once()

	err := service.RemoveContactInfo(hotelID, contactID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestListHotels(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	expectedHotels := []Hotel{
		{ID: uuid.New(), OwnerName: "John", OwnerSurname: "Doe", CompanyTitle: "Doe Ltd."},
		{ID: uuid.New(), OwnerName: "Jane", OwnerSurname: "Smith", CompanyTitle: "Smith Ltd."},
	}

	mockRepo.On("ListHotels").Return(expectedHotels, nil).Once()

	hotels, err := service.ListHotels()
	assert.NoError(t, err)
	assert.Equal(t, expectedHotels, hotels)

	mockRepo.AssertExpectations(t)
}

func TestListHotelOfficials(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	expectedOfficials := []HotelOfficial{
		{OwnerName: "John", OwnerSurname: "Doe", CompanyTitle: "Doe Ltd."},
		{OwnerName: "Jane", OwnerSurname: "Smith", CompanyTitle: "Smith Ltd."},
	}

	mockRepo.On("GetHotelOfficials").Return(expectedOfficials, nil).Once()

	officials, err := service.ListHotelOfficials()
	assert.NoError(t, err)
	assert.Equal(t, expectedOfficials, officials)

	mockRepo.AssertExpectations(t)
}

func TestGetHotelDetails(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	hotelID := uuid.New()
	expectedHotel := &Hotel{
		ID:           hotelID,
		OwnerName:    "John",
		OwnerSurname: "Doe",
		CompanyTitle: "Doe Ltd.",
	}

	mockRepo.On("GetHotelDetails", hotelID).Return(expectedHotel, nil).Once()

	hotelDetails, err := service.GetHotelDetails(hotelID)
	assert.NoError(t, err)
	assert.Equal(t, expectedHotel, hotelDetails)

	mockRepo.AssertExpectations(t)
}

func TestFetchLocationStats(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	location := "New York"
	expectedHotels := []Hotel{
		{ID: uuid.New(), ContactInfos: []ContactInfo{{InfoType: ContactTypePhone, InfoContent: "123-456"}}},
		{ID: uuid.New(), ContactInfos: []ContactInfo{{InfoType: ContactTypePhone, InfoContent: "789-101"}}},
	}
	hotelCount := len(expectedHotels)
	phoneCount := 2

	mockRepo.On("FetchHotelsByLocation", location).Return(expectedHotels, nil).Once()

	hotelCountResult, phoneCountResult, err := service.FetchLocationStats(location)
	assert.NoError(t, err)
	assert.Equal(t, hotelCount, hotelCountResult)
	assert.Equal(t, phoneCount, phoneCountResult)

	mockRepo.AssertExpectations(t)
}

func TestCreateHotel_Error(t *testing.T) {
	// Mock repository creation
	mockRepo := new(MockHotelRepository)

	// Define test hotel
	hotel := &Hotel{
		OwnerName:    "John",
		OwnerSurname: "Doe",
		CompanyTitle: "Doe Ltd.",
	}

	// Expect the Save method to be called once, but simulate an error during save
	mockRepo.On("Save", mock.MatchedBy(func(h *Hotel) bool {
		// Check if the ID is not nil and the name/surname match the test case
		return h.ID != uuid.Nil && h.OwnerName == "John" && h.OwnerSurname == "Doe"
	})).Return(fmt.Errorf("error saving hotel")).Once()

	// Create the service with the mocked repository
	service := NewService(mockRepo)

	// Call CreateHotel and assert error
	createdHotel, err := service.CreateHotel(hotel.OwnerName, hotel.OwnerSurname, hotel.CompanyTitle, nil)
	assert.Error(t, err)
	assert.Nil(t, createdHotel)

	// Ensure that the Save method was called as expected
	mockRepo.AssertExpectations(t)
}

func TestDeleteHotel_Error(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	hotelID := uuid.New()

	// Simulate an error when deleting the hotel
	mockRepo.On("Delete", hotelID).Return(fmt.Errorf("error deleting hotel")).Once()

	err := service.DeleteHotel(hotelID)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAddContactInfo_Error(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	hotelID := uuid.New()
	contact := &ContactInfo{
		InfoType:    ContactTypePhone,
		InfoContent: "123-456-7890",
	}

	// Simulate an error when adding contact info
	mockRepo.On("AddContactInfo", hotelID, contact).Return(fmt.Errorf("error adding contact info")).Once()

	err := service.AddContactInfo(hotelID, contact)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestRemoveContactInfo_Error(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	hotelID := uuid.New()
	contactID := uuid.New()

	// Simulate an error when removing contact info
	mockRepo.On("RemoveContactInfo", hotelID, contactID).Return(fmt.Errorf("error removing contact info")).Once()

	err := service.RemoveContactInfo(hotelID, contactID)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestListHotels_Empty(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	// Simulate an empty list of hotels
	mockRepo.On("ListHotels").Return([]Hotel{}, nil).Once()

	hotels, err := service.ListHotels()
	assert.NoError(t, err)
	assert.Empty(t, hotels)

	mockRepo.AssertExpectations(t)
}

func TestListHotelOfficials_Empty(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	// Simulate an empty list of hotel officials
	mockRepo.On("GetHotelOfficials").Return([]HotelOfficial{}, nil).Once()

	officials, err := service.ListHotelOfficials()
	assert.NoError(t, err)
	assert.Empty(t, officials)

	mockRepo.AssertExpectations(t)
}

func TestFetchLocationStats_ZeroHotels(t *testing.T) {
	mockRepo := new(MockHotelRepository)
	service := NewService(mockRepo)

	location := "New York"
	// Simulate zero hotels for the given location
	mockRepo.On("FetchHotelsByLocation", location).Return([]Hotel{}, nil).Once()

	hotelCount, phoneCount, err := service.FetchLocationStats(location)
	assert.NoError(t, err)
	assert.Equal(t, 0, hotelCount)
	assert.Equal(t, 0, phoneCount)

	mockRepo.AssertExpectations(t)
}
