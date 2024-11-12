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

	err = h.hotelService.AddContactInfo(hotelID, &contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contact)
}

// RemoveContactInfo handles removing contact information
func (h *Handler) RemoveContactInfo(w http.ResponseWriter, r *http.Request) {
	// Get the hotel ID and contact ID from URL parameters
	vars := mux.Vars(r)
	hotelID, err := uuid.Parse(vars["hotelID"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	contactID, err := uuid.Parse(vars["contactID"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	// Call service method to delete contact info
	err = h.hotelService.RemoveContactInfo(hotelID, contactID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content for successful deletion
	w.WriteHeader(http.StatusNoContent)
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

func (h *Handler) ListHotelOfficials(w http.ResponseWriter, r *http.Request) {
	officials, err := h.hotelService.ListHotelOfficials() // No hotelID input needed anymore
	if err != nil {
		http.Error(w, "Error retrieving hotel officials", http.StatusInternalServerError)
		return
	}

	// JSON formatında döndürme
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(officials); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetHotelDetails(w http.ResponseWriter, r *http.Request) {
	hotelID := mux.Vars(r)["hotelID"] // URL'den hotelID al
	// UUID'yi doğrulamak ve dönüştürmek için
	hotelUUID, err := uuid.Parse(hotelID)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	hotelDetails, err := h.hotelService.GetHotelDetails(hotelUUID)
	if err != nil {
		http.Error(w, "Error retrieving hotel details", http.StatusInternalServerError)
		return
	}

	// JSON formatında döndürme
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hotelDetails); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
