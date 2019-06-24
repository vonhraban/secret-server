package http

import (
	"github.com/vonhraban/secret-server/app/http/handler"

	"net/http"
)

type route struct {
	name        string
	method      string
	pattern     string
	handlerFunc http.HandlerFunc
	monitor bool
}

type routes []route

func initRoutes(secretHandler *handler.SecretHandler) routes {
	return routes{
		route{
			"ViewSecret", 
			http.MethodGet,
			"/secret/{hash}",
			secretHandler.View,
			true,
		},
		route{
			"PersistSecret",
			http.MethodPost,
			"/secret",
			secretHandler.Persist,
			true,
		},
	}	
}
