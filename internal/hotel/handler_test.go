package hotel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHotelService struct {
	mock.Mock
}

func (m *MockHotelService) CreateHotel(ownerName, ownerSurname, companyTitle string, contacts []ContactInfo) (*Hotel, error) {
	args := m.Called(ownerName, ownerSurname, companyTitle, contacts)
	return args.Get(0).(*Hotel), args.Error(1)
}

func (m *MockHotelService) DeleteHotel(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockHotelService) ListHotels() ([]Hotel, error) {
	args := m.Called()
	return args.Get(0).([]Hotel), args.Error(1)
}

func (m *MockHotelService) GetHotelDetails(id uuid.UUID) (*Hotel, error) {
	args := m.Called(id)
	return args.Get(0).(*Hotel), args.Error(1)
}

func (m *MockHotelService) AddContactInfo(hotelID uuid.UUID, contact *ContactInfo) error {
	args := m.Called(hotelID, contact)
	return args.Error(0)
}

func (m *MockHotelService) FetchLocationStats(location string) (int, int, error) {
	args := m.Called(location)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockHotelService) ListHotelOfficials() ([]HotelOfficial, error) {
	args := m.Called()
	return args.Get(0).([]HotelOfficial), args.Error(1)
}

func (m *MockHotelService) RemoveContactInfo(hotelID uuid.UUID, contactUUID uuid.UUID) error {
	args := m.Called(hotelID, contactUUID)
	return args.Error(0)
}

func TestCreateHotel_Handler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	hotel := &Hotel{
		ID:           uuid.New(),
		OwnerName:    "John",
		OwnerSurname: "Doe",
		CompanyTitle: "JD Hotels",
	}

	mockService.On("CreateHotel", "John", "Doe", "JD Hotels", mock.Anything).Return(hotel, nil)

	// Prepare the request with valid data
	requestBody := `{
		"ownerName": "John",
		"ownerSurname": "Doe",
		"companyTitle": "JD Hotels",
		"contacts": [{"phone": "123456789", "email": "contact@jd.com"}]
	}`
	req := httptest.NewRequest(http.MethodPost, "/hotels", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusCreated, rr.Code)
	var response Hotel
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, hotel.ID, response.ID)
	assert.Equal(t, "JD Hotels", response.CompanyTitle)
	mockService.AssertExpectations(t)
}

func TestDeleteHotel_Handler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	hotelID := uuid.New()

	mockService.On("DeleteHotel", hotelID).Return(nil)

	// Prepare the request with valid hotel ID
	req := httptest.NewRequest(http.MethodDelete, "/hotels/"+hotelID.String(), nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteHotel_InvalidID(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Prepare the request with invalid hotel ID
	req := httptest.NewRequest(http.MethodDelete, "/hotels/invalid-id", nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListHotels_Handler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	hotel := Hotel{
		ID:           uuid.New(),
		OwnerName:    "Alice",
		OwnerSurname: "Smith",
		CompanyTitle: "Alice's Inns",
	}

	mockService.On("ListHotels").Return([]Hotel{hotel}, nil)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/hotels", nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusOK, rr.Code)
	var hotels []Hotel
	err := json.NewDecoder(rr.Body).Decode(&hotels)
	assert.NoError(t, err)
	assert.Len(t, hotels, 1)
	assert.Equal(t, "Alice's Inns", hotels[0].CompanyTitle)
	mockService.AssertExpectations(t)
}

func TestGetHotelDetails_Handler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	hotelID := uuid.New()
	hotel := &Hotel{
		ID:           hotelID,
		OwnerName:    "Eve",
		OwnerSurname: "Brown",
		CompanyTitle: "Eve's Resorts",
	}

	mockService.On("GetHotelDetails", hotelID).Return(hotel, nil)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/hotels/"+hotelID.String(), nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusOK, rr.Code)
	var response Hotel
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, hotelID, response.ID)
	assert.Equal(t, "Eve's Resorts", response.CompanyTitle)
	mockService.AssertExpectations(t)
}

func TestGetHotelDetails_NotFound(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	hotelID := uuid.New()

	var hotel *Hotel
	// Mock the service call to return nil (hotel not found)
	mockService.On("GetHotelDetails", hotelID).Return(hotel, fmt.Errorf("hotel not found"))

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/hotels/"+hotelID.String(), nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAddContactInfo_Handler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	hotelID := uuid.New()
	contact := ContactInfo{
		InfoType:    "phone",     // Add info type
		InfoContent: "987654321", // Add info content
	}

	mockService.On("AddContactInfo", hotelID, &contact).Return(nil)

	// Prepare the request
	requestBody := `{"info_type": "phone", "info_content": "987654321"}`
	req := httptest.NewRequest(http.MethodPost, "/hotels/"+hotelID.String()+"/contacts", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestAddContactInfo_InvalidHotelID(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Prepare the request with invalid hotel ID
	requestBody := `{
		"phone": "987654321",
		"email": "contact@hotels.com"
	}`
	req := httptest.NewRequest(http.MethodPost, "/hotels/invalid-id/contacts", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRemoveContactInfo_Handler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	hotelID := uuid.New()
	contactID := uuid.New()

	mockService.On("RemoveContactInfo", hotelID, contactID).Return(nil)

	// Prepare the request
	req := httptest.NewRequest(http.MethodDelete, "/hotels/"+hotelID.String()+"/contacts/"+contactID.String(), nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockService.AssertExpectations(t)
}

func TestRemoveContactInfo_InvalidHotelID(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Prepare the request with invalid hotel ID
	contactID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/hotels/invalid-id/contacts/"+contactID.String(), nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRemoveContactInfo_InvalidContactID(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	hotelID := uuid.New()

	// Prepare the request with invalid contact ID
	req := httptest.NewRequest(http.MethodDelete, "/hotels/"+hotelID.String()+"/contacts/invalid-id", nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetHotelStats_Handler(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	location := "Paris"
	hotelCount := 10
	phoneCount := 5

	mockService.On("FetchLocationStats", location).Return(hotelCount, phoneCount, nil)

	// Prepare the request with valid location query
	req := httptest.NewRequest(http.MethodGet, "/hotels/stats?location="+location, nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusOK, rr.Code)
	var response struct {
		HotelCount int `json:"hotel_count"`
		PhoneCount int `json:"phone_count"`
	}
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, hotelCount, response.HotelCount)
	assert.Equal(t, phoneCount, response.PhoneCount)
	mockService.AssertExpectations(t)
}

func TestGetHotelStats_MissingLocation(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Prepare the request with missing location query
	req := httptest.NewRequest(http.MethodGet, "/hotels/stats", nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetHotelStats_InternalError(t *testing.T) {
	mockService := new(MockHotelService)
	handler := NewHandler(mockService)

	// Test data
	location := "Paris"

	mockService.On("FetchLocationStats", location).Return(0, 0, fmt.Errorf("internal error"))

	// Prepare the request with valid location query
	req := httptest.NewRequest(http.MethodGet, "/hotels/stats?location="+location, nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
