package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	//"github.com/gorilla/mux"
	"github.com/vonhraban/secret-server/app/http/handler"
	"github.com/vonhraban/secret-server/core/log"
	"github.com/vonhraban/secret-server/secret"
)

type service struct {
	server   *http.Server
	exitChan chan os.Signal
	logger   log.Logger
}

func New(vault secret.Vault, clock secret.Clock, logger log.Logger, port int, version string) *service {
	exitChan := make(chan os.Signal)

	secretHandler := handler.NewSecretHandler(vault, clock, logger)

	// Router
	routes := initRoutes(secretHandler)
	router := newRouter(logger, version, routes)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handlers.CORS()(router),
	}

	return &service{
		server:   server,
		exitChan: exitChan,
		logger:   logger,
	}
}

func (a *service) Serve() {
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

func (a *service) waitForExit() {
	<-a.exitChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	a.logger.Info("Server stopping...")
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error(err)
	}
}
