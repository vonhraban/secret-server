package profiler

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	//"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusProfiler struct{
	httpHitsCounter  *prometheus.CounterVec
}

func NewPrometheusProfiler() *PrometheusProfiler{
	httpHitsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
		  Namespace: "http",
		  Name: "hits",
		  Help: "Number of hits to the http endpoints",
		},
		[]string{"type"},
	)
	
	prometheus.MustRegister(httpHitsCounter)

	return &PrometheusProfiler{
		httpHitsCounter: httpHitsCounter,
	}
}

func (p *PrometheusProfiler) ServeMetrics(router *mux.Router) {
	router.Handle("/metrics", promhttp.Handler())
}

func (p *PrometheusProfiler) LogViewSecretCalled() {
	p.httpHitsCounter.WithLabelValues("view_secret").Inc()

}

func (p *PrometheusProfiler) LogPersistSecretCalled() {
	p.httpHitsCounter.WithLabelValues("persist_secret").Inc()
}