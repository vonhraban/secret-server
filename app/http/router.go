package http

import (
	"fmt"
	"net/http"
	"github.com/vonhraban/secret-server/app/http/profiler"
	"github.com/vonhraban/secret-server/core/log"

	"github.com/gorilla/mux"
)

func newRouter(logger log.Logger, version string, routes routes) *mux.Router {
	rootRouter := mux.NewRouter()

	versionRouter := rootRouter.PathPrefix(fmt.Sprintf("%s/", version)).Subrouter()
	profiler.ServePrometheus(rootRouter)

	for _, route := range routes {
		summaryVec := profiler.BuildSummaryVec(route.name, route.method + " " + route.pattern)
		versionRouter.Methods(route.method).
		Path(route.pattern).
		Name(route.name).
		Handler(profiler.WithMonitoring(route.handlerFunc, route.monitor, summaryVec))
	}

	rootRouter.NotFoundHandler = http.HandlerFunc(
		func (w http.ResponseWriter, r *http.Request) {
			logger.Warningf("unknown route requested %s%s", r.Host, r.URL.Path)
			http.Error(w, "", http.StatusNotFound)
		})

	return rootRouter
}

  