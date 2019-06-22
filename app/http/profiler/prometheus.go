package profiler

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusProfiler struct{
	viewSecretCounter prometheus.Counter
	persistSecretCounter prometheus.Counter
}

func NewPrometheusProfiler() *PrometheusProfiler{
	var (
		persistSecretCounter = promauto.NewCounter(prometheus.CounterOpts{
				Name: "http_persist_secret_calls",
				Help: "The total number of calls to persist secret endpoint",
		})
	
		viewSecretCounter = promauto.NewCounter(prometheus.CounterOpts{
			Name: "http_view_secret_calls",
			Help: "The total number of calls to view secret endpoint",
		})
	)


	return &PrometheusProfiler{
		viewSecretCounter: viewSecretCounter,
		persistSecretCounter: persistSecretCounter,
	}
}

func (p *PrometheusProfiler) ServeMetrics(router *mux.Router) {
	router.Handle("/metrics", promhttp.Handler())
}

func (p *PrometheusProfiler) LogViewSecretCalled() {
	p.viewSecretCounter.Inc()
}

func (p *PrometheusProfiler) LogPersistSecretCalled() {
	p.persistSecretCounter.Inc()
}