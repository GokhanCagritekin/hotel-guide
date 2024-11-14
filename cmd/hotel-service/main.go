package main

import (
	"hotel-guide/internal/db"
	"hotel-guide/internal/hotel"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database and ensure it closes on exit
	db.InitDB()
	defer db.CloseDB()

	// Initialize hotel repository
	hotelRepo := hotel.NewRepository()

	// Initialize hotel service
	hotelService := hotel.NewService(hotelRepo)

	// Initialize hotel handler
	hotelHandler := hotel.NewHandler(hotelService)

	// Set up router and define hotel-specific routes
	r := mux.NewRouter()
	hotelHandler.RegisterRoutes(r)

	// Start the HTTP server
	log.Println("Hotel service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
