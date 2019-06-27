package handler

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/vonhraban/secret-server/core/log"
)

func respond(w http.ResponseWriter, r *http.Request, logger log.Logger, body interface{}, responseCode int) {
	// xml if asked for specifically
	if r.Header.Get("Accept") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(body); err != nil {
			logger.Errorf("Could not encode the xml response: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
		}

		w.WriteHeader(responseCode)
		return
	}

	// assume by default json
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		logger.Errorf("Could not encode the json response: %s", err)
	}
	w.WriteHeader(responseCode)
}
