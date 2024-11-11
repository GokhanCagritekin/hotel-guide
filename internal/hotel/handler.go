package hotel

import (
	"encoding/json"
	"hotel-guide/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handler handles HTTP requests related to hotels
type Handler struct {
	hotelService *HotelService
}

// NewHandler creates a new instance of HotelHandler
func NewHandler(service *HotelService) *Handler {
	return &Handler{
		hotelService: service,
	}
}

// RegisterRoutes registers the HTTP routes for the hotel operations
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/hotels", h.CreateHotel).Methods(http.MethodPost)
	r.HandleFunc("/hotels/{id}", h.DeleteHotel).Methods(http.MethodDelete)
	r.HandleFunc("/hotels/{hotelID}/contacts", h.AddContactInfo).Methods(http.MethodPost)
	r.HandleFunc("/hotels/{hotelID}/contacts", h.RemoveContactInfo).Methods(http.MethodDelete)
	r.HandleFunc("/hotels", h.ListHotels).Methods(http.MethodGet)
}

// CreateHotel handles hotel creation
func (h *Handler) CreateHotel(w http.ResponseWriter, r *http.Request) {
	var request struct {
		OwnerName    string               `json:"ownerName"`
		OwnerSurname string               `json:"ownerSurname"`
		CompanyTitle string               `json:"companyTitle"`
		Contacts     []models.ContactInfo `json:"contacts"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hotel, err := h.hotelService.CreateHotel(request.OwnerName, request.OwnerSurname, request.CompanyTitle, request.Contacts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(hotel)
}

// DeleteHotel handles hotel deletion
func (h *Handler) DeleteHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hotelID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	err = h.hotelService.DeleteHotel(hotelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddContactInfo handles adding contact information
func (h *Handler) AddContactInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hotelID, err := uuid.Parse(vars["hotelID"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var contact models.ContactInfo
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.hotelService.AddContactInfo(hotelID, contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveContactInfo handles removing contact information
func (h *Handler) RemoveContactInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hotelID, err := uuid.Parse(vars["hotelID"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var contact models.ContactInfo
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.hotelService.RemoveContactInfo(hotelID, contact.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ListHotels handles listing all hotels
func (h *Handler) ListHotels(w http.ResponseWriter, r *http.Request) {
	hotels, err := h.hotelService.ListHotels()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hotels)
}
