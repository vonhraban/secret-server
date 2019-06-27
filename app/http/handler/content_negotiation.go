package handler

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func respond(w http.ResponseWriter, r *http.Request, body interface{}) {
	// xml if asked for specifically
	if r.Header.Get("Accept") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(body)
		return
	}

	// assume by default json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
