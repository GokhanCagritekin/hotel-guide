package main

import (
	"hotel-guide/internal/db"
	"hotel-guide/internal/hotel"
	"hotel-guide/internal/mq"
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

	// Initialize RabbitMQ
	rabbitMQ, err := mq.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close() // Close RabbitMQ connection when done

	err = rabbitMQ.InitializeQueue("reportQueue")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize services
	hotelService := hotel.NewService(hotelRepo)
	reportService := report.NewService(reportRepo, rabbitMQ) // Pass RabbitMQ to the service

	reportService.StartReportConsumer()

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
	r.HandleFunc("/reports", reportHandler.RequestReportGeneration).Methods("POST")
	r.HandleFunc("/reports", reportHandler.ListReports).Methods("GET")
	r.HandleFunc("/reports/{id}", reportHandler.GetReportByID).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
