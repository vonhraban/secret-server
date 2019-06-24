package profiler

import (
	"net/http"
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

func WithMonitoring(next http.Handler, monitor bool, summary *prometheus.SummaryVec) http.Handler {
	// Just return the next handler if route shouldn't be monitored
	if !monitor {
		return next
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		duration := time.Since(start)

		// Store duration of request in ms
		summary.WithLabelValues().Observe(float64(duration.Seconds() * 1e6))
	})
}
