package report

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReportService is the mocked version of ReportService for unit testing
type MockReportService struct {
	mock.Mock
}

// CreateReport mocks the CreateReport method
func (m *MockReportService) CreateReport(location string, hotelCount, phoneCount int) (*Report, error) {
	args := m.Called(location, hotelCount, phoneCount)
	return args.Get(0).(*Report), args.Error(1)
}

// ListReports mocks the ListReports method
func (m *MockReportService) ListReports() ([]Report, error) {
	args := m.Called()
	return args.Get(0).([]Report), args.Error(1)
}

// GetReportByID mocks the GetReportByID method
func (m *MockReportService) GetReportByID(id uuid.UUID) (*Report, error) {
	args := m.Called(id)
	return args.Get(0).(*Report), args.Error(1)
}

// RequestReportGeneration mocks the RequestReportGeneration method
func (m *MockReportService) RequestReportGeneration(location string) (*Report, error) {
	args := m.Called(location)
	return args.Get(0).(*Report), args.Error(1)
}

// UpdateReportStatus mocks the UpdateReportStatus method
func (m *MockReportService) UpdateReportStatus(id uuid.UUID, status ReportStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

// StartReportConsumer mocks the StartReportConsumer method
func (m *MockReportService) StartReportConsumer() {
	m.Called()
}

// fetchLocationStats mocks the fetchLocationStats method
func (m *MockReportService) fetchLocationStats(location string) (int, int, error) {
	args := m.Called(location)
	return args.Int(0), args.Int(1), args.Error(2)
}

// Test RequestReportGeneration
func TestRequestReportGeneration_Handler(t *testing.T) {
	mockService := new(MockReportService)
	handler := NewHandler(mockService)

	// Test data
	reportID := uuid.New()
	report := &Report{
		ID:       reportID,
		Location: "Paris",
		Status:   Pending,
	}

	mockService.On("RequestReportGeneration", "Paris").Return(report, nil)

	// Prepare the request
	requestBody := `{"location": "Paris"}`
	req := httptest.NewRequest(http.MethodPost, "/reports", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusCreated, rr.Code)
	var response Report
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, reportID, response.ID)
	assert.Equal(t, "Paris", response.Location)
	assert.Equal(t, Pending, response.Status)
	mockService.AssertExpectations(t)
}

// Test ListReports
func TestListReports_Handler(t *testing.T) {
	mockService := new(MockReportService)
	handler := NewHandler(mockService)

	// Test data
	report := Report{
		ID:       uuid.New(),
		Location: "Paris",
		Status:   "Completed",
	}

	mockService.On("ListReports").Return([]Report{report}, nil)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/reports", nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusOK, rr.Code)
	var reports []Report
	err := json.NewDecoder(rr.Body).Decode(&reports)
	assert.NoError(t, err)
	assert.Len(t, reports, 1)
	assert.Equal(t, "Paris", reports[0].Location)
	mockService.AssertExpectations(t)
}

// Test GetReportByID
func TestGetReportByID_Handler(t *testing.T) {
	mockService := new(MockReportService)
	handler := NewHandler(mockService)

	// Test data
	reportID := uuid.New()
	report := &Report{
		ID:       reportID,
		Location: "Paris",
		Status:   "Completed",
	}

	mockService.On("GetReportByID", reportID).Return(report, nil)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/reports/"+reportID.String(), nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusOK, rr.Code)
	var response Report
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, reportID, response.ID)
	assert.Equal(t, "Paris", response.Location)
	assert.Equal(t, Completed, response.Status)
	mockService.AssertExpectations(t)
}

// Test GetReportByID_NotFound
func TestGetReportByID_NotFound(t *testing.T) {
	mockService := new(MockReportService)
	handler := NewHandler(mockService)

	// Test data
	reportID := uuid.New()
	var report *Report
	// Mock the service call to return nil, nil (no report found)
	mockService.On("GetReportByID", reportID).Return(report, nil)

	// Prepare the request
	req := httptest.NewRequest(http.MethodGet, "/reports/"+reportID.String(), nil)
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockService.AssertExpectations(t)
}

// Test RequestReportGeneration_InvalidRequestBody
func TestRequestReportGeneration_InvalidRequestBody(t *testing.T) {
	mockService := new(MockReportService)
	handler := NewHandler(mockService)

	// Prepare the request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/reports", bytes.NewBufferString(`{invalid json}`))
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRequestReportGeneration_EmptyLocation(t *testing.T) {
	mockService := new(MockReportService)
	handler := NewHandler(mockService)

	// Prepare the request with empty location
	req := httptest.NewRequest(http.MethodPost, "/reports", bytes.NewBufferString(`{"location": ""}`))
	rr := httptest.NewRecorder()

	// Register routes and handle request
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.ServeHTTP(rr, req)

	// Assert status code and response body
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
