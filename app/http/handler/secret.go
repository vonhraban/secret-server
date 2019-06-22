package handler

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/vonhraban/secret-server/secret"
	"github.com/vonhraban/secret-server/secret/cmd"
	"github.com/vonhraban/secret-server/secret/query"
	"github.com/vonhraban/secret-server/core/log"
	"github.com/vonhraban/secret-server/app/http/profiler"
)

type secretHandler struct {
	vault secret.Vault
	clock secret.Clock
	logger log.Logger
	profiler *profiler.PrometheusProfiler
}

func NewSecretHandler(
		vault secret.Vault,
		clock secret.Clock,
		logger log.Logger,
		profiler *profiler.PrometheusProfiler,
	) *secretHandler {
	return &secretHandler{
		vault: vault,
		clock: clock,
		logger: logger,
		profiler: profiler,
	}
}

func (h *secretHandler) Persist(w http.ResponseWriter, r *http.Request) {
	h.profiler.LogViewSecretCalled()
	request, err := buildPersistSecretRequestFromHTTPRequest(r)
	if err != nil {
		switch err.(type) {
		case *EmptyValueError:
			h.logger.Warningf("Validation error %s", err)
			response := NewErrorResponse(err.Error())

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(response)

		default:
			h.logger.Error(err)
			response := NewErrorResponse("Internal Error")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)		
			
		}

		return
	}

	hash := uuid.NewV4().String()

	command := cmd.NewAddSecretCommand(
		h.vault,
		h.clock,
		hash,
		request.secret,
		request.expireAfterViews,
		request.expireAfter,
	)

	if err := command.Execute(); err != nil {
			h.logger.Error(err)
			response := NewErrorResponse("Internal Error")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)		
			
			return
	}

	q := query.NewGetSecretQuery(h.vault, hash)
	storedSecret, err := q.Execute()
	if err != nil {
		h.logger.Error(err)
		response := NewErrorResponse("Internal Error")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)		

		return
	}

	response := buildPersistSecretResponseFromSecret(*storedSecret)

	// xml if asked for specifically
	if r.Header.Get("Accept") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
		return
	}

	// assume by default json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

func (h *secretHandler) View(w http.ResponseWriter, r *http.Request) {
	h.profiler.LogPersistSecretCalled()
	params := mux.Vars(r)
	hash := params["hash"]

	q := query.NewGetSecretQuery(h.vault, hash)

	storedSecret, err := q.Execute()
	if err != nil {
		// TODO! Catch specificallt not found error
		http.Error(w, "", http.StatusNotFound)
		return
	}

	response := buildViewSecretResponseFromSecret(*storedSecret)

	decreaseViewsCmd := cmd.NewDecreaseRemainingViewsCommand(h.vault, hash)
	if err := decreaseViewsCmd.Execute(); err != nil {
		panic(err)
	}

	// Since we now decreased the number of available views in a store secret, we need to descrease it in the response too
	// and we want to do it only if there are more that 0 remaining views
	if response.RemainingViews > 0 {
		response.RemainingViews--
	}

	// TODO! Remoce duplication - perhaps middleware?
	// xml if asked for specifically
	if r.Header.Get("Accept") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
		return
	}

	// assume by default json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}
