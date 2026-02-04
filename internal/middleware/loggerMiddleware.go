package middleware

import (
	"agrigation_api/pkg/logger/logger"
	"context"
	"fmt"
	"net/http"
	"time"
)

// LoggerMiddleware - middleware для логгов
func LoggerMiddleware(logs *logger.Log, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs.Info(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v",
			r.RemoteAddr, r.URL, r.Method, time.Now()), logger.GetPlace())

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "logger", logs)))
	})
}
