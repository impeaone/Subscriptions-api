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
	t, err := time.Parse("01-2006", dateStr)
	fmt.Println(t)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(strings.TrimSpace(id))
}
