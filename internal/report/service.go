package report

import (
	"fmt"
	"hotel-guide/models"

	"github.com/google/uuid"
)

// Service provides operations on reports
type Service struct {
	reportRepo ReportRepository
}

// NewService creates a new instance of ReportService
func NewService(repo ReportRepository) *Service {
	return &Service{
		reportRepo: repo,
	}
}

// CreateReport creates a new report
func (s *Service) CreateReport(location string, hotelCount, phoneCount int) (*models.Report, error) {
	report := models.NewReport(location, hotelCount, phoneCount)
	err := s.reportRepo.CreateReport(report)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// UpdateReportStatus updates the status of a report
func (s *Service) UpdateReportStatus(reportID uuid.UUID, status models.ReportStatus) error {
	err := s.reportRepo.UpdateReportStatus(reportID, status)
	if err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}
	return nil
}

// ListReports lists all reports
func (s *Service) ListReports() ([]models.Report, error) {
	return s.reportRepo.ListReports()
}

// GetReportByID gets a report by its ID
func (s *Service) GetReportByID(id uuid.UUID) (*models.Report, error) {
	return s.reportRepo.GetReportByID(id)
}
