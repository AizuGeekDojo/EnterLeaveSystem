package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AizuGeekDojo/EnterLeaveSystem/server/db"
	"github.com/AizuGeekDojo/EnterLeaveSystem/server/handler"
	"github.com/AizuGeekDojo/EnterLeaveSystem/server/utils"
	"golang.org/x/net/websocket"
)

const (
	// ServerPort is the port the HTTP server listens on
	ServerPort = ":3000"
	// ShutdownTimeout is the maximum time to wait for graceful shutdown
	ShutdownTimeout = 30 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	// Setup signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Open database connection
	log.Println("Opening database...")
	d, err := db.OpenDB()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		if err := d.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	h := handler.NewHandler(d)

	// Initialize Slack Socket Mode client
	log.Println("Initializing Slack Socket Mode...")
	if err := utils.InitSlackSocketMode(); err != nil {
		log.Printf("Warning: Failed to initialize Slack Socket Mode: %v", err)
	}

	// Setup HTTP handlers
	log.Println("Setting up HTTP handlers...")
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("dist"))))
	mux.Handle("/socket/readCard", websocket.Handler(h.ReadCardHandler))
	mux.HandleFunc("/api/user", h.UserAPIHandler)
	mux.HandleFunc("/api/log", h.LogAPIHandler)

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:              ServerPort,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start NFC card reader in background
	log.Println("Starting NFC card reader...")
	go handler.ReadCard(d)

	// Start cron scheduler
	log.Println("Starting cron scheduler...")
	if err := utils.CronInit(d); err != nil {
		return fmt.Errorf("failed to initialize cron: %w", err)
	}

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Server starting on %s", ServerPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Println("Shutdown signal received, starting graceful shutdown...")
	}

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("Server shutdown completed successfully")
	return nil
}
