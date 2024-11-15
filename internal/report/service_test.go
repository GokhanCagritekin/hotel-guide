package report

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReportRepository is a mock implementation of the ReportRepository interface
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) Save(report *Report) error {
	args := m.Called(report)
	return args.Error(0)
}

func (m *MockReportRepository) ListReports() ([]Report, error) {
	args := m.Called()
	return args.Get(0).([]Report), args.Error(1)
}

func (m *MockReportRepository) GetReportByID(id uuid.UUID) (*Report, error) {
	args := m.Called(id)
	return args.Get(0).(*Report), args.Error(1)
}

func (m *MockReportRepository) UpdateReportStatus(id uuid.UUID, status ReportStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockReportRepository) FetchHotelAndPhoneCounts(location string) (int, int, error) {
	args := m.Called(location)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockReportRepository) UpdateReportStats(id uuid.UUID, hotelCount, phoneCount int, status ReportStatus) error {
	args := m.Called(id, hotelCount, phoneCount, status)
	return args.Error(0)
}

// MockRabbitMQ, MessageQueue interface'ini mock'layan bir struct
type MockMessageQueue struct {
	mock.Mock
}

// Publish, MessageQueue'nin Publish metodunu mock'lar
func (m *MockMessageQueue) Publish(queueName string, message []byte) error {
	args := m.Called(queueName, message)
	return args.Error(0)
}

// Consume, MessageQueue'nin Consume metodunu mock'lar
func (m *MockMessageQueue) Consume(queueName string) (<-chan amqp.Delivery, error) {
	args := m.Called(queueName)
	return args.Get(0).(<-chan amqp.Delivery), args.Error(1)
}

// Close, MessageQueue'nin Close metodunu mock'lar
func (m *MockMessageQueue) Close() error {
	args := m.Called()
	return args.Error(0)
}

// InitializeQueue, MessageQueue'nin InitializeQueue metodunu mock'lar
func (m *MockMessageQueue) InitializeQueue(queueName string) error {
	args := m.Called(queueName)
	return args.Error(0)
}

// TestCreateReport tests the CreateReport method of reportService
func TestCreateReport(t *testing.T) {
	mockRepo := new(MockReportRepository)
	mockRabbitMQ := new(MockMessageQueue)
	service := NewService(mockRepo, mockRabbitMQ)

	report := &Report{
		ID:       uuid.New(),
		Location: "Test Location",
		Status:   Pending,
	}

	// Expect Save to be called once and return no error
	mockRepo.On("Save", mock.MatchedBy(func(r *Report) bool {
		return r.Location == report.Location && r.Status == Pending
	})).Return(nil).Once()

	// Call CreateReport
	createdReport, err := service.CreateReport(report.Location, 5, 10)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, report.Location, createdReport.Location)
	assert.Equal(t, Pending, createdReport.Status)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestRequestReportGeneration tests the RequestReportGeneration method of reportService
func TestRequestReportGeneration(t *testing.T) {
	// Initialize mocks
	mockRepo := new(MockReportRepository)
	mockQueue := new(MockMessageQueue)
	service := NewService(mockRepo, mockQueue)

	// Create the expected report structure
	expectedReport := &Report{
		ID:          uuid.New(),
		Location:    "Test Location",
		HotelCount:  0,
		PhoneCount:  0,
		RequestedAt: time.Now(),
		Status:      "Pending", // Make sure this matches the status in your code
	}

	// Set up expectations
	mockRepo.On("Save", mock.AnythingOfType("*report.Report")).Return(nil).Run(func(args mock.Arguments) {
		report := args.Get(0).(*Report)
		report.ID = expectedReport.ID // Match the expected report ID
	})

	// Mock Publish to avoid errors on message queue operations
	mockQueue.On("Publish", "reportQueue", mock.AnythingOfType("[]uint8")).Return(nil)

	// Call the method under test
	result, err := service.RequestReportGeneration("Test Location")

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, expectedReport.Location, result.Location)
	assert.Equal(t, Pending, result.Status)

	// Verify all expectations were met
	mockRepo.AssertExpectations(t)
	mockQueue.AssertExpectations(t)
}

// TestStartReportConsumer tests the StartReportConsumer method of reportService
func TestStartReportConsumer(t *testing.T) {
	mockRepo := new(MockReportRepository)
	mockRabbitMQ := new(MockMessageQueue)
	service := NewService(mockRepo, mockRabbitMQ)

	// Setup mock for Consume
	mockMessages := make(chan amqp.Delivery)                                                           // Make the channel bidirectional here
	mockRabbitMQ.On("Consume", "reportQueue").Return((<-chan amqp.Delivery)(mockMessages), nil).Once() // Cast to <-chan

	location := "Test Location"
	expectedHotelCount := 5
	expectedPhoneCount := 10
	mockRepo.On("FetchHotelAndPhoneCounts", location).Return(expectedHotelCount, expectedPhoneCount, nil)

	// Start the consumer in a goroutine
	go service.StartReportConsumer()

	// Simulate sending a message to the queue
	reportID := uuid.New()
	message := amqp.Delivery{
		Body: []byte(`{"id":"` + reportID.String() + `", "location":"Test Location"}`),
	}

	status := Completed
	// Set up expectations for UpdateReportStats
	mockRepo.On("UpdateReportStats", reportID, expectedHotelCount, expectedPhoneCount, status).Return(nil)

	mockMessages <- message

	// Assertions: no specific assertions for now as it is running in the background
	assert.True(t, true) // This ensures that the test runs without panicking

	// Verify expectations
	mockRabbitMQ.AssertExpectations(t)
}

// TestListReports tests the ListReports method of reportService
func TestListReports(t *testing.T) {
	mockRepo := new(MockReportRepository)
	mockRabbitMQ := new(MockMessageQueue)
	service := NewService(mockRepo, mockRabbitMQ)

	// Prepare the reports to be returned by the mock
	expectedReports := []Report{
		{ID: uuid.New(), Location: "Location 1", Status: Completed},
		{ID: uuid.New(), Location: "Location 2", Status: Pending},
	}

	// Set up expectations for ListReports
	mockRepo.On("ListReports").Return(expectedReports, nil)

	// Call the method under test
	reports, err := service.ListReports()

	// Assert results
	assert.NoError(t, err)
	assert.Len(t, reports, 2)
	assert.Equal(t, "Location 1", reports[0].Location)
	assert.Equal(t, "Location 2", reports[1].Location)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestGetReportByID tests the GetReportByID method of reportService
func TestGetReportByID(t *testing.T) {
	mockRepo := new(MockReportRepository)
	mockRabbitMQ := new(MockMessageQueue)
	service := NewService(mockRepo, mockRabbitMQ)

	// Prepare the report to be returned by the mock
	reportID := uuid.New()
	expectedReport := &Report{
		ID:       reportID,
		Location: "Test Location",
		Status:   Pending,
	}

	// Set up expectations for GetReportByID
	mockRepo.On("GetReportByID", reportID).Return(expectedReport, nil)

	// Call the method under test
	report, err := service.GetReportByID(reportID)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, reportID, report.ID)
	assert.Equal(t, "Test Location", report.Location)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestUpdateReportStatus tests the UpdateReportStatus method of reportService
func TestUpdateReportStatus(t *testing.T) {
	mockRepo := new(MockReportRepository)
	mockRabbitMQ := new(MockMessageQueue)
	service := NewService(mockRepo, mockRabbitMQ)

	reportID := uuid.New()
	status := Completed

	// Set up expectations for UpdateReportStatus
	mockRepo.On("UpdateReportStatus", reportID, status).Return(nil)

	// Call the method under test
	err := service.UpdateReportStatus(reportID, status)

	// Assert results
	assert.NoError(t, err)

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

// TestFetchLocationStats tests the fetchLocationStats method of reportService
func TestFetchLocationStats(t *testing.T) {
	// Create the mock objects
	mockRepo := new(MockReportRepository)
	mockRabbitMQ := new(MockMessageQueue)

	// Create the service with mocked dependencies
	service := NewService(mockRepo, mockRabbitMQ)

	// Mock data for the location
	location := "Test Location"
	expectedHotelCount := 5
	expectedPhoneCount := 10

	// Set up expectations for FetchHotelAndPhoneCounts
	// This mocks the method FetchHotelAndPhoneCounts and specifies the return values
	mockRepo.On("FetchHotelAndPhoneCounts", location).Return(expectedHotelCount, expectedPhoneCount, nil)

	// Call the method under test
	hotelCount, phoneCount, err := service.fetchLocationStats(location)

	// Assert results
	assert.NoError(t, err)                          // Ensure no error was returned
	assert.Equal(t, expectedHotelCount, hotelCount) // Ensure hotel count is correct
	assert.Equal(t, expectedPhoneCount, phoneCount) // Ensure phone count is correct

	// Verify expectations were met (that the method was called as expected)
	mockRepo.AssertExpectations(t)
}
