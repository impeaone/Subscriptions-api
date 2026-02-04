package tools

import (
	"os"
	"strconv"
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
