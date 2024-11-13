package report

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// ReportHandler struct to handle HTTP requests
type ReportHandler struct {
	reportService ReportService
}

// NewReportHandler creates a new ReportHandler
func NewHandler(service ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: service,
	}
}

// RegisterRoutes registers report-related routes
func (h *ReportHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/reports", h.ListReports).Methods(http.MethodGet)
	r.HandleFunc("/reports/{id}", h.GetReportByID).Methods(http.MethodGet)
	r.HandleFunc("/reports", h.RequestReportGeneration).Methods(http.MethodPost)
}

// RequestReportGeneration handles the creation of a new report
func (h *ReportHandler) RequestReportGeneration(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Location string `json:"location"`
	}

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate location
	if req.Location == "" {
		http.Error(w, "Location must not be empty", http.StatusBadRequest)
		return
	}

	// Call the service to request a new report generation
	report, err := h.reportService.RequestReportGeneration(req.Location)
	if err != nil {
		log.Error().Err(err).Msg("Error creating report")
		http.Error(w, fmt.Sprintf("Error creating report: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the created report
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(report); err != nil {
		log.Error().Err(err).Msg("Error encoding response")
	}
}

// ListReports handles fetching all reports
func (h *ReportHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.reportService.ListReports()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching reports")
		http.Error(w, fmt.Sprintf("Error fetching reports: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the list of reports
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reports); err != nil {
		log.Error().Err(err).Msg("Error encoding response")
	}
}

// GetReportByID handles fetching a specific report by ID
func (h *ReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request) {
	// Get report ID from URL params
	params := mux.Vars(r)
	id, err := uuid.Parse(params["id"])
	if err != nil {
		http.Error(w, "Invalid report ID format", http.StatusBadRequest)
		return
	}

	// Fetch the report by ID
	report, err := h.reportService.GetReportByID(id)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching report")
		http.Error(w, fmt.Sprintf("Error fetching report: %v", err), http.StatusInternalServerError)
		return
	}
	if report == nil {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}

	// Return the report
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(report); err != nil {
		log.Error().Err(err).Msg("Error encoding response")
	}
}
