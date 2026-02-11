package middleware

import (
	"encoding/json"
	"net/http"
)

// ShutdownMiddleware - middleware проверяющий закрыт ли канал, для graceful shutdown
func ShutdownMiddleware(exitChan chan struct{}, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-exitChan:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":   "service_unavailable",
				"message": "Service is shutting down",
			})
			return
		default:
			next.ServeHTTP(w, r)
		}
	})
}
