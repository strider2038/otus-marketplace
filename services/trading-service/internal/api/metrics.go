package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	RequestCount     *prometheus.CounterVec
	RequestDurations *prometheus.HistogramVec
}

func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		RequestCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "request_count",
				Help:      "Request count",
			},
			[]string{"method", "route", "http_status"},
		),
		RequestDurations: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "request_latency_seconds",
				Help:      "Request latency in seconds",
			},
			[]string{"method", "route"},
		),
	}
}

func MetricsMiddleware(next http.Handler, metrics *Metrics) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		route := mux.CurrentRoute(request)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(metrics.RequestDurations.WithLabelValues(request.Method, path))
		defer timer.ObserveDuration()

		sw := &statusWriter{writer: writer}
		next.ServeHTTP(sw, request)
		status := sw.status

		metrics.RequestCount.WithLabelValues(request.Method, path, strconv.Itoa(status)).Inc()
	})
}

type statusWriter struct {
	writer http.ResponseWriter
	status int
}

func (w *statusWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *statusWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.status = statusCode

	w.writer.WriteHeader(statusCode)
}
