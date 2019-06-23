package http

import (
	"github.com/vonhraban/secret-server/app/http/handler"

	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Monitor bool
}

type Routes []Route

func initRoutes(secretHandler *handler.SecretHandler) Routes {
	return Routes{
		Route{
			"ViewSecret", 
			http.MethodGet,
			"/secret/{hash}",
			secretHandler.View,
			true,
		},
		Route{
			"PersistSecret",
			http.MethodPost,
			"/secret",
			secretHandler.Persist,
			true,
		},
	}	
}
