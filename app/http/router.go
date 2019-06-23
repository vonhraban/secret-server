package http

import (
	"net/http"
	"time"
	//"strconv"
	"fmt"
	//"github.com/vonhraban/secret-server/app/http/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"


	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

func buildSummaryVec(metricName string, metricHelp string) *prometheus.SummaryVec {
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

func withMonitoring(next http.Handler, route Route, summary *prometheus.SummaryVec) http.Handler {
	// Just return the next handler if route shouldn't be monitored
	if !route.Monitor {
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

// NewRouter returns root router and version router
func NewRouter(version string, routes Routes) *mux.Router {
	rootRouter := mux.NewRouter()
	rootRouter.Handle("/metrics", promhttp.Handler())

	versionRouter := rootRouter.PathPrefix(fmt.Sprintf("%s/", version)).Subrouter()

	for _, route := range routes {
		// create summaryVec for endpoint
		summaryVec := buildSummaryVec(route.Name, route.Method + " " + route.Pattern)
		versionRouter.Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(withMonitoring(route.HandlerFunc, route, summaryVec)) // <-- CHAINING HERE!!!
	}
	//logrus.Infoln("Successfully initialized routes including Prometheus.")
	return rootRouter
}
