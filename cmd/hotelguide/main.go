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
	// Initialize the database and ensure it closes on exit
	db.InitDB()
	defer db.CloseDB()

	// Initialize repositories
	hotelRepo := hotel.NewRepository()
	reportRepo := report.NewRepository()

	// Retrieve RabbitMQ connection URL from the mq package
	rabbitMQURL, err := mq.NewRabbitMQURL()
	if err != nil {
		log.Fatalf("Failed to get RabbitMQ URL: %v", err)
	}

	// Initialize RabbitMQ connection and queue setup
	rabbitMQ, err := mq.NewRabbitMQ(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Initialize services
	hotelService := hotel.NewService(hotelRepo)
	reportService := report.NewService(reportRepo, rabbitMQ)

	// Start the report consumer for processing asynchronous tasks
	reportService.StartReportConsumer()

	// Initialize handlers
	hotelHandler := hotel.NewHandler(hotelService)
	reportHandler := report.NewHandler(reportService)

	// Set up router and define routes
	r := mux.NewRouter()
	hotelHandler.RegisterRoutes(r)
	reportHandler.RegisterRoutes(r)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))
}
