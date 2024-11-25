package handlers

import (
	"context"
	"net/http"
	"time"

	"zim-iot-data-api/database"

	"github.com/sirupsen/logrus"
)

// Probes holds dependencies for probe handlers
type Probes struct {
	DB     *database.DB
	Logger *logrus.Logger
}

// NewProbes creates a new Probes instance
func NewProbes(db *database.DB, logger *logrus.Logger) *Probes {
	return &Probes{
		DB:     db,
		Logger: logger,
	}
}

// ReadinessProbe handles the /readiness endpoint
func (p *Probes) ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := p.DB.PingContext(ctx); err != nil {
		p.Logger.WithError(err).Warn("Readiness probe failed: Database not ready")
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// LivenessProbe handles the /liveness endpoint
func (p *Probes) LivenessProbe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
