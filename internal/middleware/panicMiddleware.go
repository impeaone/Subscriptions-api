package middleware

import (
	logger2 "agrigation_api/pkg/logger"
	"agrigation_api/pkg/logger/logger"
	"net/http"
)

func PanicMiddleware(log logger2.MyLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic: "+err.(error).Error(), logger.GetPlace())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
