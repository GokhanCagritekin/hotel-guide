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
	// DB bağlantısını başlatma
	db.InitDB()
	defer db.CloseDB()

	// Repo ve service'leri oluşturma
	hotelRepo := hotel.NewRepository()
	reportRepo := report.NewRepository()

	hotelService := hotel.NewService(hotelRepo)
	reportService := report.NewService(reportRepo)

	// Create the handler instances
	hotelHandler := hotel.NewHandler(hotelService)
	reportHandler := report.NewHandler(reportService)

	// Router oluşturma
	r := mux.NewRouter()

	// API route'ları
	r.HandleFunc("/hotels", hotelHandler.CreateHotel).Methods("POST")
	r.HandleFunc("/hotels/{id}", hotelHandler.DeleteHotel).Methods("DELETE")
	r.HandleFunc("/hotels", hotelHandler.ListHotels).Methods("GET")
	r.HandleFunc("/hotels/{hotelID}/contacts", hotelHandler.AddContactInfo).Methods("POST")
	r.HandleFunc("/hotels/{hotelID}/contacts/{contactID}", hotelHandler.RemoveContactInfo).Methods("DELETE")
	r.HandleFunc("/hotels/officials", hotelHandler.ListHotelOfficials).Methods("GET")
	r.HandleFunc("/hotels/{hotelID}", hotelHandler.GetHotelDetails).Methods("GET")
	r.HandleFunc("/reports", reportHandler.CreateReport).Methods("POST")
	r.HandleFunc("/reports/{id}", reportHandler.GetReportByID).Methods("GET")

	// Sunucuyu başlatma
	log.Fatal(http.ListenAndServe(":8080", r))
}
