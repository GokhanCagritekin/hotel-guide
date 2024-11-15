package main

import (
	"hotel-guide/internal/db"
	"hotel-guide/internal/mq"
	"hotel-guide/internal/report"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database and ensure it closes on exit
	dbInstance, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.CloseDB(dbInstance)

	if err := dbInstance.AutoMigrate(&report.Report{}); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Initialize report repository
	reportRepo := report.NewRepository(dbInstance)

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

	err = rabbitMQ.InitializeQueue("reportQueue")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize report service with RabbitMQ dependency
	reportService := report.NewService(reportRepo, rabbitMQ)

	// Start the report consumer for processing asynchronous tasks
	reportService.StartReportConsumer()

	// Initialize report handler
	reportHandler := report.NewHandler(reportService)

	// Set up router and define report-specific routes
	r := mux.NewRouter()
	reportHandler.RegisterRoutes(r)

	// Start the HTTP server
	log.Println("Report service is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}
