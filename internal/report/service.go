package report

import (
	"encoding/json"
	"fmt"
	"hotel-guide/internal/mq"
	"hotel-guide/models"
	"log"

	"github.com/google/uuid"
)

// ReportService interface defines the methods for report-related operations
type ReportService interface {
	CreateReport(location string, hotelCount, phoneCount int) (*models.Report, error)
	ListReports() ([]models.Report, error)
	GetReportByID(id uuid.UUID) (*models.Report, error)
	RequestReportGeneration(location string) (*models.Report, error)
	UpdateReportStatus(id uuid.UUID, status models.ReportStatus) error
	StartReportConsumer()
	fetchLocationStats(location string) (int, int, error)
}

// reportService struct implements the ReportService interface
type reportService struct {
	reportRepo ReportRepository
	rabbitMQ   *mq.RabbitMQ
}

// NewReportService creates a new instance of reportService
func NewService(repo ReportRepository, rabbitMQ *mq.RabbitMQ) ReportService {
	return &reportService{
		reportRepo: repo,
		rabbitMQ:   rabbitMQ,
	}
}

// CreateReport creates a new report with the provided details
func (s *reportService) CreateReport(location string, hotelCount, phoneCount int) (*models.Report, error) {
	report := models.NewReport(location, hotelCount, phoneCount)
	err := s.reportRepo.Save(report)
	if err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}
	return report, nil
}

// ListReports retrieves a list of all reports
func (s *reportService) ListReports() ([]models.Report, error) {
	return s.reportRepo.ListReports()
}

// GetReportByID retrieves the details of a report by its ID
func (s *reportService) GetReportByID(id uuid.UUID) (*models.Report, error) {
	return s.reportRepo.GetReportByID(id)
}

// RequestReportGeneration handles the creation of a new report and sends it to the RabbitMQ queue
func (s *reportService) RequestReportGeneration(location string) (*models.Report, error) {
	// Create a new report with "Pending" status
	report := models.NewReport(location, 0, 0) // Initial counts set to 0
	report.Status = models.Pending
	err := s.reportRepo.Save(report)
	if err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	// Marshal the report ID and location to JSON
	reportRequest := struct {
		ID       uuid.UUID `json:"id"`
		Location string    `json:"location"`
	}{
		ID:       report.ID,
		Location: location,
	}

	reportJSON, err := json.Marshal(reportRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal report request to JSON: %w", err)
	}

	// Send the JSON-formatted report request to RabbitMQ
	err = s.rabbitMQ.Publish("reportQueue", reportJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to publish report generation request: %w", err)
	}

	return report, nil
}

// UpdateReportStatus updates the status of an existing report
func (s *reportService) UpdateReportStatus(id uuid.UUID, status models.ReportStatus) error {
	return s.reportRepo.UpdateReportStatus(id, status)
}

// StartReportConsumer consumes messages from RabbitMQ and processes reports
func (s *reportService) StartReportConsumer() {
	messages, err := s.rabbitMQ.Consume("reportQueue")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start consumer: %w", err))
	}

	go func() {
		for msg := range messages {
			var request struct {
				ID       uuid.UUID `json:"id"`
				Location string    `json:"location"`
			}
			err := json.Unmarshal(msg.Body, &request)
			if err != nil {
				log.Printf("Invalid report request in message: %v", err)
				continue
			}

			// Fetch hotel and phone counts for the specified location
			hotelCount, phoneCount, err := s.fetchLocationStats(request.Location)
			if err != nil {
				log.Printf("Failed to fetch location stats for %s: %v", request.Location, err)
				continue
			}

			// Update the report with the fetched stats and set status to Completed
			err = s.reportRepo.UpdateReportStats(request.ID, hotelCount, phoneCount, models.Completed)
			if err != nil {
				log.Printf("Failed to update report status for report ID %s: %v", request.ID, err)
				continue
			}

			log.Printf("Report %s has been successfully processed with %d hotels and %d phones", request.ID, hotelCount, phoneCount)
		}
	}()
}

// fetchLocationStats fetches hotel and phone counts for a given location.
func (s *reportService) fetchLocationStats(location string) (int, int, error) {
	hotelCount, phoneCount, err := s.reportRepo.FetchHotelAndPhoneCounts(location)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch hotel and phone counts for location %s: %w", location, err)
	}
	return hotelCount, phoneCount, nil
}
