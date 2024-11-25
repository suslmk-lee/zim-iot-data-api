package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zim-iot-data-api/config"
	"zim-iot-data-api/database"
	"zim-iot-data-api/handlers"
	"zim-iot-data-api/utils"
)

func main() {
	// Initialize logger
	logger := utils.InitLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	// If not production, load from config.properties
	if cfg.Profile != "prod" {
		err = config.LoadConfigFromIni(cfg)
		if err != nil {
			logger.WithError(err).Fatal("Failed to load configuration from config.properties")
		}
	}

	// Initialize database
	db, err := database.NewDB(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}

	// Initialize handlers
	iotHandlers := handlers.NewIoTHandlers(db, logger)
	probes := handlers.NewProbes(db, logger)

	// Define HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/iot-data", iotHandlers.GetIoTData)
	mux.HandleFunc("/iot-data/latest", iotHandlers.GetLatestIoTData)
	mux.HandleFunc("/readiness", probes.ReadinessProbe)
	mux.HandleFunc("/liveness", probes.LivenessProbe)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting server on port %s...", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	// Close the database connection
	if err := db.Close(); err != nil {
		logger.WithError(err).Fatal("Failed to close database")
	}

	logger.Info("Server exiting")
}
