package middleware

import (
	"net/http"
	"time"

	"bisnis-rinzi/packages/backend/logger"
)

// LoggerMiddleware mencatat metode HTTP, rute yang diakses, status respon, dan durasi eksekusi request.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Membungkus ResponseWriter untuk menangkap status code HTTP
		rw := &responseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		logger.Info("[%s] %s %s - Status: %d (%s)",
			r.RemoteAddr, r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterInterceptor) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
