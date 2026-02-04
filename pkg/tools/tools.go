package tools

import (
	"agrigation_api/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnv - считывает переменную окружения, если нет, возвращает дефолт-значение
func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvAsBool считывает значение переменной окружения как булево или возвращает значение по умолчанию,
// если переменная не установлена или не может быть преобразована в булево
func GetEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := GetEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// GetEnvAsInt считывает значение переменной окружения как int или возвращает значение по умолчанию,
// если переменная не установлена или не может быть преобразована в int
func GetEnvAsInt(key string, defaultValue int) int {
	if valueStr := GetEnv(key, ""); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// WriteJSON - хелпер для JSON ответов
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError - хелпер для ошибок
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, models.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// ParseMonthYear : парсим дату "MM-YYYY" -> time.Time
func ParseMonthYear(dateStr string) (time.Time, error) {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	monthStr := parts[0]
	yearStr := parts[1]

	// Определяем формат: "07-2025" или "2025-07"
	var month, year int
	var err error

	if len(monthStr) <= 2 && len(yearStr) == 4 {
		// Формат "MM-YYYY"
		month, err = strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			return time.Time{}, fmt.Errorf("invalid month: %s", monthStr)
		}
		year, err = strconv.Atoi(yearStr)
		if err != nil || year < 2000 || year > 2100 {
			return time.Time{}, fmt.Errorf("invalid year: %s", yearStr)
		}
	} else if len(monthStr) == 4 && len(yearStr) <= 2 {
		// Формат "YYYY-MM"
		year, err = strconv.Atoi(monthStr)
		if err != nil || year < 2000 || year > 2100 {
			return time.Time{}, fmt.Errorf("invalid year: %s", monthStr)
		}
		month, err = strconv.Atoi(yearStr)
		if err != nil || month < 1 || month > 12 {
			return time.Time{}, fmt.Errorf("invalid month: %s", yearStr)
		}
	} else {
		return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(strings.TrimSpace(id))
}
