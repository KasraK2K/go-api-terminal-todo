package main

import (
	"context"     // Provides context for request cancellation
	"errors"      // Error handling utilities
	"log"         // Logging utility
	"net/http"    // HTTP server functionalities
	"os"          // OS-level functions
	"os/signal"   // Captures system signals (e.g., SIGTERM)
	"runtime"     // Provides system information like CPU cores
	"strconv"     // Converts numbers to strings
	"sync"        // Synchronization utilities (WaitGroup)
	"syscall"     // Low-level OS signal handling (e.g., SIGTERM)
	"time"        // Time utilities (for delays)
	"todo/routes" // Custom package for setting up HTTP routes
)

func startServer(port string, wg *sync.WaitGroup, errChan chan error, shutdownChan <-chan struct{}) {
	defer wg.Done() // Ensure the WaitGroup counter is decremented when the function exits

	mux := routes.SetupRoutes() // Set up HTTP handlers
	server := &http.Server{
		Addr:    ":" + port, // Bind server to the specified port
		Handler: mux,        // Use the HTTP router from `routes`
	}

	// Channel to handle server shutdown
	idleConnsClosed := make(chan struct{})

	go func() {
		<-shutdownChan // Wait for shutdown signal
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Server on port %s shutdown error: %v", port, err)
		}
		close(idleConnsClosed) // Close channel when shutdown is complete
	}()

	log.Printf("Server running on port %s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- err // Send error to the error channel if server fails
	}

	<-idleConnsClosed // Wait for shutdown before exiting function
}

func bootstrap() {
	numWorkers := runtime.NumCPU()      // Get the number of CPU cores
	ports := make([]string, numWorkers) // Create a list of ports

	for i := 0; i < numWorkers; i++ {
		ports[i] = strconv.Itoa(3000 + i) // Generate ports 3000, 3001, etc.
	}

	var wg sync.WaitGroup
	errChan := make(chan error, numWorkers) // Buffered error channel
	shutdownChan := make(chan struct{})     // Channel to signal shutdown

	// Start all servers
	for _, port := range ports {
		wg.Add(1)                                        // Increment WaitGroup counter
		go startServer(port, &wg, errChan, shutdownChan) // Start server in a goroutine
	}

	// OS signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Log once all servers have started
	go func() {
		time.Sleep(10 * time.Millisecond)
		log.Println("All servers started successfully")
	}()

	// Monitor errors
	go func() {
		for err := range errChan {
			log.Printf("Server error: %v", err)
		}
	}()

	// Block until an OS signal is received
	<-signalChan
	log.Println("Shutdown signal received. Stopping all servers...")

	// Trigger shutdown for all servers
	close(shutdownChan)

	// Wait for all servers to shut down before exiting
	wg.Wait()
	log.Println("All servers shut down. Exiting.")
}
