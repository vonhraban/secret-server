package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/vonhraban/secret-server/app/http/handler"
	"github.com/vonhraban/secret-server/secret"
	"github.com/vonhraban/secret-server/core/log"
)

type Service struct {
	server   *http.Server
	exitChan chan os.Signal
	logger log.Logger
}

func New(vault secret.Vault, clock secret.Clock, logger log.Logger) *Service {
	exitChan := make(chan os.Signal)
	router := mux.NewRouter()

	v1 := router.PathPrefix("/v1").Subrouter()

	secretHandler := handler.NewSecretHandler(vault, clock, logger)

	v1.HandleFunc("/secret", secretHandler.Persist).Methods(http.MethodPost)
	v1.HandleFunc("/secret/{hash}", secretHandler.View).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    ":80",
		Handler: handlers.CORS()(router),
	}

	return &Service{
		server:   server,
		exitChan: exitChan,
		logger: logger,
	}
}

func (a *Service) Serve() {
	go func() {
		err := a.server.ListenAndServe()
		if err != http.ErrServerClosed {
			a.logger.Fatal(err)
		}
		a.logger.Info("Server shut down")
	}()

	a.logger.Info("Server started")

	signal.Notify(a.exitChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	a.waitForExit()
}

func (a *Service) waitForExit() {
	<-a.exitChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	a.logger.Info("Server stopping...")
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error(err)
	}
}
