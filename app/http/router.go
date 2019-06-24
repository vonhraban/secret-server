package http

import (
	"fmt"
	"github.com/vonhraban/secret-server/app/http/profiler"

	"github.com/gorilla/mux"
)

func newRouter(version string, routes routes) *mux.Router {
	rootRouter := mux.NewRouter()

	versionRouter := rootRouter.PathPrefix(fmt.Sprintf("%s/", version)).Subrouter()
	profiler.ServePrometheus(rootRouter)

	for _, route := range routes {
		summaryVec := profiler.BuildSummaryVec(route.name, route.method + " " + route.pattern)
		versionRouter.Methods(route.method).
		Path(route.pattern).
		Name(route.name).
		Handler(profiler.WithMonitoring(route.handlerFunc, route.monitor, summaryVec)) // <-- CHAINING HERE!!!
	}
	//logrus.Infoln("Successfully initialized routes including Prometheus.")
	return rootRouter
}
