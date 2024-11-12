package main

import (
	"hotel-guide/internal/db"
	"hotel-guide/internal/hotel"
	"hotel-guide/internal/report"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize DB and close on exit
	db.InitDB()
	defer db.CloseDB()

	// Initialize repositories
	hotelRepo := hotel.NewRepository()
	reportRepo := report.NewRepository()

	// Initialize services
	hotelService := hotel.NewService(hotelRepo)
	reportService := report.NewService(reportRepo)

	// Initialize handlers
	hotelHandler := hotel.NewHandler(hotelService)
	reportHandler := report.NewHandler(reportService)

	// Initialize router and define routes
	r := mux.NewRouter()

	r.HandleFunc("/hotels", hotelHandler.CreateHotel).Methods("POST")
	r.HandleFunc("/hotels/{id}", hotelHandler.DeleteHotel).Methods("DELETE")
	r.HandleFunc("/hotels", hotelHandler.ListHotels).Methods("GET")
	r.HandleFunc("/hotels/{hotelID}/contacts", hotelHandler.AddContactInfo).Methods("POST")
	r.HandleFunc("/hotels/{hotelID}/contacts/{contactID}", hotelHandler.RemoveContactInfo).Methods("DELETE")
	r.HandleFunc("/hotels/officials", hotelHandler.ListHotelOfficials).Methods("GET")
	r.HandleFunc("/hotels/{hotelID}", hotelHandler.GetHotelDetails).Methods("GET")
	r.HandleFunc("/reports", reportHandler.CreateReport).Methods("POST")
	r.HandleFunc("/reports/{id}", reportHandler.GetReportByID).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
