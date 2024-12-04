package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"zim-iot-data-api/database"
	"zim-iot-data-api/models"

	"github.com/sirupsen/logrus"
)

// IoTHandlers holds dependencies for IoT-related handlers
type IoTHandlers struct {
	DB     *database.DB
	Logger *logrus.Logger
}

// NewIoTHandlers creates a new IoTHandlers instance
func NewIoTHandlers(db *database.DB, logger *logrus.Logger) *IoTHandlers {
	return &IoTHandlers{
		DB:     db,
		Logger: logger,
	}
}

// getIoTData handles the /iot-data endpoint
func (h *IoTHandlers) GetIoTData(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	// Input Validation for recent_count
	recentCountStr := r.URL.Query().Get("recent_count")
	recentCount, err := strconv.Atoi(recentCountStr)
	if err != nil || recentCount <= 0 || recentCount > 1000 { // 최대 1000으로 제한
		recentCount = 10
	}

	// Input Validation for within_one_hour
	withinOneHourStr := r.URL.Query().Get("within_one_hour")
	withinOneHour, err := strconv.ParseBool(withinOneHourStr)
	if err != nil {
		withinOneHour = false
	}

	var data []models.IoTData
	if withinOneHour {
		oneHourAgo := time.Now().Add(-1 * time.Hour)
		query := `SELECT device, timestamp, pro_ver, minor_ver, sn, model, tyield, dyield, pf, pmax, pac, sac, uab, ubc, uca, ia, ib, ic, freq, tmod, tamb, mode, qac, bus_capacitance, ac_capacitance, pdc, pmax_lim, smax_lim, is_sent, reg_timestamp 
				  FROM iot_data 
				  WHERE timestamp >= $1 
				  ORDER BY timestamp DESC 
				  LIMIT $2`
		data, err = h.queryIoTData(ctx, query, oneHourAgo, recentCount)
	} else {
		query := `SELECT device, timestamp, pro_ver, minor_ver, sn, model, tyield, dyield, pf, pmax, pac, sac, uab, ubc, uca, ia, ib, ic, freq, tmod, tamb, mode, qac, bus_capacitance, ac_capacitance, pdc, pmax_lim, smax_lim, is_sent, reg_timestamp 
				  FROM iot_data 
				  ORDER BY timestamp DESC 
				  LIMIT $1`
		data, err = h.queryIoTData(ctx, query, recentCount)
	}

	if err != nil {
		h.Logger.WithError(err).Error("Failed to query IoT data")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to encode JSON")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetLatestIoTData handles the /iot-data/latest endpoint
func (h *IoTHandlers) GetLatestIoTData(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	query := `SELECT device, timestamp, pro_ver, minor_ver, sn, model, tyield, dyield, pf, pmax, pac, sac, uab, ubc, uca, ia, ib, ic, freq, tmod, tamb, mode, qac, bus_capacitance, ac_capacitance, pdc, pmax_lim, smax_lim, is_sent, reg_timestamp 
			  FROM iot_data 
			  ORDER BY timestamp DESC 
			  LIMIT 1`
	data, err := h.queryIoTData(ctx, query)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to query latest IoT data")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "No data found", http.StatusNotFound)
		return
	}

	response := models.SingleIoTDataResponse{Data: data[0]}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Logger.WithError(err).Error("Failed to encode JSON")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// queryIoTData queries the database for IoT data using context
func (h *IoTHandlers) queryIoTData(ctx context.Context, query string, args ...interface{}) ([]models.IoTData, error) {
	rows, err := h.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.IoTData
	for rows.Next() {
		var d models.IoTData
		err = rows.Scan(
			&d.Device,
			&d.Timestamp,
			&d.ProVer,
			&d.MinorVer,
			&d.SN,
			&d.Model,
			&d.TYield,
			&d.DYield,
			&d.PF,
			&d.PMax,
			&d.PAC,
			&d.SAC,
			&d.UAB,
			&d.UBC,
			&d.UCA,
			&d.IA,
			&d.IB,
			&d.IC,
			&d.Freq,
			&d.TMod,
			&d.TAmb,
			&d.Mode,
			&d.QAC,
			&d.BusCapacitance,
			&d.ACCapacitance,
			&d.PDC,
			&d.PMaxLim,
			&d.SMaxLim,
			&d.IsSent,
			&d.RegTimestamp,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}
