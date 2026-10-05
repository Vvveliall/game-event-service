package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPErrorsTotal      *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	EventsProcessedTotal prometheus.Counter
}

func New() *Metrics {
	m := &Metrics{
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "game_event_http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),
		HTTPErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "game_event_http_errors_total",
				Help: "Total number of HTTP requests that returned an error status.",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "game_event_http_request_duration_seconds",
				Help: "HTTP request duration in seconds.",
			},
			[]string{"method", "path"},
		),
		EventsProcessedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "game_event_events_processed_total",
				Help: "Total number of events processed by workers.",
			},
		),
	}

	prometheus.MustRegister(m.HTTPRequestsTotal)
	prometheus.MustRegister(m.HTTPErrorsTotal)
	prometheus.MustRegister(m.HTTPRequestDuration)
	prometheus.MustRegister(m.EventsProcessedTotal)

	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}

func (m *Metrics) RecordHTTPRequest(
	method string,
	path string,
	status int,
	duration time.Duration,
) {
	statusCode := strconv.Itoa(status)

	m.HTTPRequestsTotal.WithLabelValues(
		method,
		path,
		statusCode,
	).Inc()

	if status >= http.StatusInternalServerError {
		m.HTTPErrorsTotal.WithLabelValues(
			method,
			path,
			statusCode,
		).Inc()
	}

	m.HTTPRequestDuration.WithLabelValues(
		method,
		path,
	).Observe(duration.Seconds())
}

func (m *Metrics) RecordEventProcessed() {
	m.EventsProcessedTotal.Inc()
}
