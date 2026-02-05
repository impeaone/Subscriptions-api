package handlers

import (
	"agrigation_api/pkg/logger/logger"
	"agrigation_api/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HealthCheck godoc
// @Summary Health check
// @Description Check if API is running
// @Tags system
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "agrigation-api",
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	if errEncode := json.NewEncoder(w).Encode(response); errEncode != nil {
		h.logs.Error(fmt.Sprintf("Json-Encode Error: %v", errEncode), logger.GetPlace())
	}
}
