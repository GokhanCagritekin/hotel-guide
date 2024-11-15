package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRepository defines report database operations
type ReportRepository interface {
	Save(report *Report) error
	ListReports() ([]Report, error)
	GetReportByID(id uuid.UUID) (*Report, error)
	UpdateReportStatus(id uuid.UUID, status ReportStatus) error
	UpdateReportStats(reportID uuid.UUID, hotelCount, phoneCount int, status ReportStatus) error
	FetchHotelAndPhoneCounts(location string) (int, int, error)
}
type reportRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// Save saves a new report
func (r *reportRepository) Save(report *Report) error {
	return r.db.Create(report).Error
}

// ListReports lists all reports
func (r *reportRepository) ListReports() ([]Report, error) {
	var reports []Report
	err := r.db.Find(&reports).Error
	return reports, err
}

// GetReportByID fetches a report by its ID
func (r *reportRepository) GetReportByID(id uuid.UUID) (*Report, error) {
	var report Report
	err := r.db.First(&report, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

// UpdateReportStatus updates the status of a report
func (r *reportRepository) UpdateReportStatus(id uuid.UUID, status ReportStatus) error {
	return r.db.Model(&Report{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateReportStats updates the hotel count, phone count, and status of a report
func (r *reportRepository) UpdateReportStats(reportID uuid.UUID, hotelCount, phoneCount int, status ReportStatus) error {
	return r.db.Model(&Report{}).
		Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"hotel_count": hotelCount,
			"phone_count": phoneCount,
			"status":      status,
		}).Error
}

// FetchHotelAndPhoneCounts fetches hotel and phone counts by location from hotel-service
func (r *reportRepository) FetchHotelAndPhoneCounts(location string) (int, int, error) {
	var hotelServiceURL = os.Getenv("HOTEL_SERVICE_URL")
	location = url.QueryEscape(location)
	url := fmt.Sprintf("%s/hotels/stats?location=%s", hotelServiceURL, location)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch hotel and phone counts from hotel-service: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		HotelCount int `json:"hotel_count"`
		PhoneCount int `json:"phone_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		return 0, 0, fmt.Errorf("failed to decode hotel and phone counts response: %w", err)
	}
	return result.HotelCount, result.PhoneCount, nil
}
