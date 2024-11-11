package report

import (
	"fmt"

	"github.com/google/uuid"
)

// Service provides operations on reports
type Service struct {
	reportRepo ReportRepository
}

// ReportRepository defines the methods required to interact with the data storage
type ReportRepository interface {
	CreateReport(report *Report) error
	UpdateReportStatus(reportID uuid.UUID, status ReportStatus) error
	ListReports() ([]Report, error)
	GetReportByID(id uuid.UUID) (*Report, error)
}

// NewService creates a new instance of ReportService
func NewService(repo ReportRepository) *Service {
	return &Service{
		reportRepo: repo,
	}
}

// CreateReport creates a new report
func (s *Service) CreateReport(location string, hotelCount, phoneCount int) (*Report, error) {
	report := NewReport(location, hotelCount, phoneCount)
	err := s.reportRepo.CreateReport(report)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// UpdateReportStatus updates the status of a report
func (s *Service) UpdateReportStatus(reportID uuid.UUID, status ReportStatus) error {
	err := s.reportRepo.UpdateReportStatus(reportID, status)
	if err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}
	return nil
}

// ListReports lists all reports
func (s *Service) ListReports() ([]Report, error) {
	return s.reportRepo.ListReports()
}

// GetReportByID gets a report by its ID
func (s *Service) GetReportByID(id uuid.UUID) (*Report, error) {
	return s.reportRepo.GetReportByID(id)
}
