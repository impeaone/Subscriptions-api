package models

import "time"

// HealthResponse структура ответа для health check
type HealthResponse struct {
	Status    string    `json:"status" example:"ok"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00Z"`
	Service   string    `json:"service" example:"user-api"`
	Version   string    `json:"version" example:"1.0.0"`
}
