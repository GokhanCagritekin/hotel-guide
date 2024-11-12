package report

import (
	"errors"
	"time"

	"hotel-guide/internal/db"
	"hotel-guide/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRepository provides an interface for report operations
type ReportRepository interface {
	CreateReport(report *models.Report) error
	GetReportByID(id uuid.UUID) (*models.Report, error)
	ListReports() ([]models.Report, error)
	UpdateReportStatus(id uuid.UUID, status models.ReportStatus) error
}

// reportRepository is the implementation of the ReportRepository interface
type reportRepository struct {
	db *gorm.DB
}

// NewRepository creates a new reportRepository instance
func NewRepository() ReportRepository {
	return &reportRepository{
		db: db.DB,
	}
}

// CreateReport creates a new report record
func (r *reportRepository) CreateReport(report *models.Report) error {
	report.ID = generateUUID() // Generates a new UUID for the report
	report.RequestedAt = time.Now()
	report.Status = "Preparing"
	return r.db.Create(report).Error
}

// GetReportByID fetches a report by its UUID
func (r *reportRepository) GetReportByID(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	if err := r.db.First(&report, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if the record is not found
		}
		return nil, err
	}
	return &report, nil
}

// ListReports fetches all reports
func (r *reportRepository) ListReports() ([]models.Report, error) {
	var reports []models.Report
	if err := r.db.Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

// UpdateReportStatus updates the status of a report by its UUID
func (r *reportRepository) UpdateReportStatus(id uuid.UUID, status models.ReportStatus) error {
	return r.db.Model(&models.Report{}).Where("id = ?", id).Update("status", status).Error
}

// generateUUID generates a new UUID
func generateUUID() uuid.UUID {
	return uuid.New()
}
