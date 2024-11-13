package report

import (
	"errors"
	"fmt"

	"hotel-guide/internal/db"
	"hotel-guide/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRepository defines report database operations
type ReportRepository interface {
	Save(report *models.Report) error
	ListReports() ([]models.Report, error)
	GetReportByID(id uuid.UUID) (*models.Report, error)
	UpdateReportStatus(id uuid.UUID, status models.ReportStatus) error
	GetHotelStatsByLocation(location string, hotels *[]models.Hotel) error
	UpdateReportStats(reportID uuid.UUID, hotelCount, phoneCount int, status models.ReportStatus) error
}

// reportRepository implements ReportRepository
type reportRepository struct {
	db *gorm.DB
}

// NewRepository returns a new reportRepository instance
func NewRepository() ReportRepository {
	return &reportRepository{db: db.DB}
}

// Save saves a new report
func (r *reportRepository) Save(report *models.Report) error {
	return r.db.Create(report).Error
}

// ListReports lists all reports
func (r *reportRepository) ListReports() ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Find(&reports).Error
	return reports, err
}

// GetReportByID fetches a report by its ID
func (r *reportRepository) GetReportByID(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := r.db.First(&report, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

// UpdateReportStatus updates the status of a report
func (r *reportRepository) UpdateReportStatus(id uuid.UUID, status models.ReportStatus) error {
	return r.db.Model(&models.Report{}).Where("id = ?", id).Update("status", status).Error
}

// GetHotelStatsByLocation fetches hotel stats by location
func (r *reportRepository) GetHotelStatsByLocation(location string, hotels *[]models.Hotel) error {
	err := r.db.Joins("JOIN contact_infos ON contact_infos.hotel_id = hotels.id").
		Where("contact_infos.info_type IN (?)", []string{"location", "phone"}).
		Where("contact_infos.info_content = ?", location).
		Preload("ContactInfos"). // Preload the ContactInfos relation
		Find(hotels).Error

	if err != nil {
		return fmt.Errorf("failed to get hotels by location: %w", err)
	}

	return nil
}

// UpdateReportStats updates the hotel count, phone count, and status of a report
func (r *reportRepository) UpdateReportStats(reportID uuid.UUID, hotelCount, phoneCount int, status models.ReportStatus) error {
	return r.db.Model(&models.Report{}).
		Where("id = ?", reportID).
		Updates(map[string]interface{}{
			"hotel_count": hotelCount,
			"phone_count": phoneCount,
			"status":      status,
		}).Error
}
