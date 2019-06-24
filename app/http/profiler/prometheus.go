package profiler

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

func BuildSummaryVec(metricName string, metricHelp string) *prometheus.SummaryVec {
	summaryVec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "http",
			Name:      metricName,
			Help:      metricHelp,
			Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.005, 0.99: 0.001},

		},
		[]string{},
	)
	prometheus.Register(summaryVec)
	return summaryVec
}

func ServePrometheus(router *mux.Router) {
	router.Handle("/metrics", promhttp.Handler())
}