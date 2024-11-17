package main

import (
	"context"
	"hotel-guide/internal/db"
	"hotel-guide/internal/hotel"
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
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.CloseDB(dbInstance)

	// Run migrations
	if err := dbInstance.AutoMigrate(&hotel.Hotel{}, &hotel.ContactInfo{}); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Initialize hotel repository
	hotelRepo := hotel.NewRepository(dbInstance)

	// Initialize hotel service
	hotelService := hotel.NewService(hotelRepo)

	// Initialize hotel handler
	hotelHandler := hotel.NewHandler(hotelService)

	// Set up router and define hotel-specific routes
	r := mux.NewRouter()
	hotelHandler.RegisterRoutes(r)

	// Setup HTTP server with graceful shutdown capabilities
	server := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	// Run the server in a goroutine so that we can listen for shutdown signals
	go func() {
		log.Println("Hotel service is running on port 8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Graceful shutdown logic
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal to gracefully shutdown the server
	<-stop
	log.Println("Shutting down the hotel service...")

	// Define a graceful shutdown timeout (e.g., 5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Hotel service stopped gracefully")
}
