package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "chaosboard_http_requests_total",
		Help: "Total Http requests",
	},
		[]string{"method", "path", "status"},
	)

	HttpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "chaosboard_http_requests_duration",
		Help:    "Duration of http requests",
		Buckets: prometheus.DefBuckets,
	},
		[]string{"method", "path"},
	)

	ExperimentsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "chaosboard_experiments_total",
		Help: "Total no. of experiments",
	},
		[]string{"type", "status"},
	)

	ExperimentsRunning = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "chaosboard_experiments_active",
		Help: "Total no. of experiments currently running",
	})
)

type responseWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func TrackRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK, false}
		next.ServeHTTP(rw, r)

		HttpRequests.WithLabelValues(
			r.Method,
			r.URL.Path,
			fmt.Sprintf("%d", rw.status),
		).Inc()

		HttpDuration.WithLabelValues(r.Method, r.URL.Path).Observe(float64(time.Since(start).Seconds()))
	})
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.written {
		return
	}
	rw.status = code
	rw.written = true
	rw.ResponseWriter.WriteHeader(code)
}
