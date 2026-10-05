package middleware

import (
	"net/http"
	"time"

	"game-event-service/internal/metrics"
)

func Metrics(m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			writer := &responseWriter{
				ResponseWriter: w,
			}

			next.ServeHTTP(writer, r)

			status := writer.statusCode
			if status == 0 {
				status = http.StatusOK
			}

			m.RecordHTTPRequest(
				r.Method,
				r.URL.Path,
				status,
				time.Since(start),
			)
		})
	}
}
