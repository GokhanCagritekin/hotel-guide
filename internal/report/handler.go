package report

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler handles HTTP requests related to reports
type Handler struct {
	reportService *Service
}

// NewHandler creates a new instance of ReportHandler
func NewHandler(service *Service) *Handler {
	return &Handler{
		reportService: service,
	}
}

// RegisterRoutes registers the HTTP routes for the report operations
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/reports", h.CreateReport).Methods(http.MethodPost)
	r.HandleFunc("/reports/{id}", h.GetReportByID).Methods(http.MethodGet)
	r.HandleFunc("/reports", h.ListReports).Methods(http.MethodGet)
}

// CreateReport handles report creation
func (h *Handler) CreateReport(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Location   string `json:"location"`
		HotelCount int    `json:"hotelCount"`
		PhoneCount int    `json:"phoneCount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	report, err := h.reportService.CreateReport(request.Location, request.HotelCount, request.PhoneCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

// ListReports handles listing all reports
func (h *Handler) ListReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.reportService.ListReports()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reports)
}

// GetReportByID handles getting a report by its ID
func (h *Handler) GetReportByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	report, err := h.reportService.GetReportByID(reportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}
