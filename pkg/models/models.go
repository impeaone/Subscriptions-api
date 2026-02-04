package models

import (
	"github.com/google/uuid"
	"time"
)

// HealthResponse структура ответа для health check
type HealthResponse struct {
	Status    string    `json:"status" example:"ok"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00Z"`
	Service   string    `json:"service" example:"user-api"`
	Version   string    `json:"version" example:"1.0.0"`
}

// Subscription - подписка пользователя
type Subscription struct {
	UserID      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	StartDate   string    `json:"start_date"`         // "07-2025"
	EndDate     *string   `json:"end_date,omitempty"` // "12-2025" или null
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrUpdateRequest - запрос на создание/обновление
type CreateOrUpdateRequest struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date,omitempty"`
}

// DeleteRequest - запрос на удаление
type DeleteRequest struct {
	ServiceName string    `json:"service_name"`
	UserID      uuid.UUID `json:"user_id"`
}

// ErrorResponse - ошибка API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// ListResponse - ответ со списком
type ListResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Total         int            `json:"total"`
}

// UserSubscriptionsResponse - подписки пользователя
type UserSubscriptionsResponse struct {
	UserID        uuid.UUID      `json:"user_id"`
	Subscriptions []Subscription `json:"subscriptions"`
	MonthlyTotal  int            `json:"monthly_total"`
	Currency      string         `json:"currency"`
}

// CalculateTotalRequest - запрос для подсчета суммы
type CalculateTotalRequest struct {
	UserID      uuid.UUID `json:"user_id,omitempty"`      // опционально
	ServiceName string    `json:"service_name,omitempty"` // опционально
	StartMonth  string    `json:"start_month"`            // "01-2024" начало периода
	EndMonth    string    `json:"end_month"`              // "12-2024" конец периода
}

// CalculateTotalResponse - ответ для подсчета суммы
type CalculateTotalResponse struct {
	Success  bool       `json:"success"`
	Total    int        `json:"total"`
	Currency string     `json:"currency"`
	Period   PeriodInfo `json:"period"`
	Filters  FilterInfo `json:"filters"`
}

type PeriodInfo struct {
	StartMonth string `json:"start_month"`
	EndMonth   string `json:"end_month"`
}

type FilterInfo struct {
	UserID      *string `json:"user_id,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
}
