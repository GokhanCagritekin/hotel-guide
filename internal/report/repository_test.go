package report

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSaveReport(t *testing.T) {
	// Mock setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	// Mock SQLite version query
	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).
		WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// GORM setup
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Repository setup
	repo := NewRepository(gormDB)

	// Mock data
	expectedTime := time.Now()
	report := &Report{
		ID:          uuid.New(),
		Location:    "Test Location",
		HotelCount:  0,
		PhoneCount:  0,
		RequestedAt: expectedTime,
		Status:      Pending,
	}

	// Mock expectations
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `reports`").
		WithArgs(
			report.Location,
			report.HotelCount,
			report.PhoneCount,
			expectedTime,
			report.Status,
			report.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test Save method
	err = repo.Save(report)
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestListReports_Repository(t *testing.T) {
	// Mock database setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database connection: %v", err)
	}
	defer db.Close()

	// Mock SQLite version query
	mock.ExpectQuery(`(?i)^SELECT sqlite_version\(\)$`).
		WillReturnRows(sqlmock.NewRows([]string{"sqlite_version"}).AddRow("3.32.3"))

	// Initialize GORM DB
	gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %v", err)
	}

	// Create repository instance
	repo := NewRepository(gormDB)

	// Expected results
	mockReports := []Report{
		{
			ID:         uuid.New(),
			Location:   "Location 1",
			HotelCount: 5,
			PhoneCount: 10,
			Status:     Completed,
		},
		{
			ID:         uuid.New(),
			Location:   "Location 2",
			HotelCount: 3,
			PhoneCount: 7,
			Status:     Pending,
		},
	}

	// Expectations for ListReports
	rows := sqlmock.NewRows([]string{"id", "location", "hotel_count", "phone_count", "status"}).
		AddRow(mockReports[0].ID.String(), mockReports[0].Location, mockReports[0].HotelCount, mockReports[0].PhoneCount, mockReports[0].Status).
		AddRow(mockReports[1].ID.String(), mockReports[1].Location, mockReports[1].HotelCount, mockReports[1].PhoneCount, mockReports[1].Status)

	mock.ExpectQuery(`SELECT \* FROM ` + "`reports`").
		WillReturnRows(rows)

	// Call ListReports method
	reports, err := repo.ListReports()
	assert.NoError(t, err)
	assert.Equal(t, len(mockReports), len(reports))

	// Validate results
	for i, report := range reports {
		assert.Equal(t, mockReports[i].ID, report.ID)
		assert.Equal(t, mockReports[i].Location, report.Location)
		assert.Equal(t, mockReports[i].HotelCount, report.HotelCount)
		assert.Equal(t, mockReports[i].PhoneCount, report.PhoneCount)
		assert.Equal(t, mockReports[i].Status, report.Status)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestFetchHotelAndPhoneCounts(t *testing.T) {
	// Set up environment variable for hotel service URL
	hotelServiceURL := "http://localhost:8080" // Ensure it points to localhost for mock server
	os.Setenv("HOTEL_SERVICE_URL", hotelServiceURL)
	defer os.Unsetenv("HOTEL_SERVICE_URL")

	// Expected response
	mockLocation := "Test Location"
	mockHotelCount := 5
	mockPhoneCount := 10

	// Start a mock HTTP server using httptest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, fmt.Sprintf("/hotels/stats?location=%s", url.QueryEscape(mockLocation)), r.URL.String())
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"hotel_count": %d, "phone_count": %d}`, mockHotelCount, mockPhoneCount)
	}))
	defer server.Close()

	// Ensure environment variable points to the mock server
	os.Setenv("HOTEL_SERVICE_URL", server.URL)

	// Initialize repository with a dummy DB (not used in this test)
	gormDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) // Using in-memory SQLite for simplicity
	repo := NewRepository(gormDB)

	// Call FetchHotelAndPhoneCounts
	hotelCount, phoneCount, err := repo.FetchHotelAndPhoneCounts(mockLocation)
	assert.NoError(t, err)
	assert.Equal(t, mockHotelCount, hotelCount)
	assert.Equal(t, mockPhoneCount, phoneCount)
}
