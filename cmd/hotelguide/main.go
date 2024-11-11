package main

import (
	"hotel-guide/internal/hotel"
	"hotel-guide/internal/report"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Dependency Injection
	hotelService := hotel.NewService(nil)
	reportService := report.NewService(nil)

	hotelHandler := hotel.NewHandler(hotelService)
	reportHandler := report.NewHandler(reportService)

	r := mux.NewRouter()
	hotelHandler.RegisterRoutes(r)
	reportHandler.RegisterRoutes(r)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
