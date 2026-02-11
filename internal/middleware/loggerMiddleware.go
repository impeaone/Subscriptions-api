package middleware

import (
	logger2 "agrigation_api/pkg/logger"
	"agrigation_api/pkg/logger/logger"
	"context"
	"fmt"
	"net/http"
)

// LoggerMiddleware - middleware для логгов
func LoggerMiddleware(logs logger2.MyLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs.Info(fmt.Sprintf("Client: %s; EndPoint: %s; Method: %s; Time: %v",
			r.RemoteAddr, r.URL, r.Method, logger.TimeFormat), logger.GetPlace())

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "logger", logs)))
	})
}
