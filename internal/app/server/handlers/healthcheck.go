package handlers

import (
	"agrigation_api/pkg/logger/logger"
	"agrigation_api/pkg/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("logger").(*logger.Log)

	response := models.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "agrigation-api",
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	if errEncode := json.NewEncoder(w).Encode(response); errEncode != nil {
		log.Error(fmt.Sprintf("Json-Encode Error: %v", errEncode), logger.GetPlace())
	}
}
