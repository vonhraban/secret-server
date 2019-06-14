package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	// router
	router := mux.NewRouter()

	v1 := router.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/", helloWorldHandler).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    ":80",
		Handler: handlers.CORS()(router),
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
		fmt.Println("Server shut down")
	}()

	log.Println("Server started")

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("Server stopping...")
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
