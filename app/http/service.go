package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/vonhraban/secret-server/app/http/handler"
	"github.com/vonhraban/secret-server/secret"
)

type Service struct {
	server   *http.Server
	exitChan chan os.Signal
}

func New(vault secret.Vault, clock secret.Clock) *Service {
	exitChan := make(chan os.Signal)
	router := mux.NewRouter()

	v1 := router.PathPrefix("/v1").Subrouter()

	secretHandler := handler.NewSecretHandler(vault, clock)

	v1.HandleFunc("/secret", secretHandler.Persist).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    ":80",
		Handler: handlers.CORS()(router),
	}

	return &Service{
		server:   server,
		exitChan: exitChan,
	}
}

func (a *Service) Serve() {
	go func() {
		err := a.server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
		fmt.Println("Server shut down")
	}()

	log.Println("Server started")

	signal.Notify(a.exitChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	a.waitForExit()
}

func (a *Service) waitForExit() {
	<-a.exitChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("Server stopping...")
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
