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
	logger := utils.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	db, err := database.NewDB(cfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}

	iotHandlers := handlers.NewIoTHandlers(db, logger)
	probes := handlers.NewProbes(db, logger)

	mux := http.NewServeMux()

	// 각 핸들러에 CORS 미들웨어 적용
	mux.Handle("/iot-data", utils.CORSMiddleware(http.HandlerFunc(iotHandlers.GetIoTData)))
	mux.Handle("/iot-data/latest", utils.CORSMiddleware(http.HandlerFunc(iotHandlers.GetLatestIoTData)))
	mux.Handle("/readiness", utils.CORSMiddleware(http.HandlerFunc(probes.ReadinessProbe)))
	mux.Handle("/liveness", utils.CORSMiddleware(http.HandlerFunc(probes.LivenessProbe)))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	go func() {
		logger.Infof("Starting server on port %s...", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	if err := db.Close(); err != nil {
		logger.WithError(err).Fatal("Failed to close database")
	}

	logger.Info("Server exiting")
}
