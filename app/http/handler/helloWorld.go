package handler

import (
	"encoding/json"
	"net/http"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func HelloNameHandler(w http.ResponseWriter, r *http.Request) {

	request := struct {
		Name string `json:"name"`
	}{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "err", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Hello " + request.Name))
}
