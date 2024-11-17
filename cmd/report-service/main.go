package main

import (
	"context"
	"hotel-guide/internal/db"
	"hotel-guide/internal/mq"
	"hotel-guide/internal/report"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database and ensure it closes on exit
	dbInstance, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer db.CloseDB(dbInstance)

	// Run migrations
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

	if err := rabbitMQ.InitializeQueue("reportQueue"); err != nil {
		log.Fatalf("Error initializing RabbitMQ queue: %v", err)
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

	// Setup HTTP server with graceful shutdown capabilities
	server := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	// Run the server in a goroutine so that we can listen for shutdown signals
	go func() {
		log.Println("Report service is running on port 8082")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Graceful shutdown logic
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal to gracefully shutdown the server
	<-stop
	log.Println("Shutting down the report service...")

	// Define a graceful shutdown timeout (e.g., 5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Report service stopped gracefully")
}
