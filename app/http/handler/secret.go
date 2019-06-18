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
)

type secretHandler struct {
	vault secret.Vault
	clock secret.Clock
}

func NewSecretHandler(vault secret.Vault, clock secret.Clock) *secretHandler {
	return &secretHandler{
		vault: vault,
		clock: clock,
	}
}

func (h *secretHandler) Persist(w http.ResponseWriter, r *http.Request) {
	//panic(fmt.Sprintf("%+v", r))
	request, err := persistSecretRequestFromHTTPRequest(r)
	if err != nil {
		panic(err)
	}
	// TODO! Validation

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
		panic(err)
	}

	q := query.NewGetSecretQuery(h.vault, hash)
	storedSecret, err := q.Execute()
	if err != nil {
		panic(err)
	}

	response := persistSecretResponseFromSecret(*storedSecret)

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
	// TODO! Validation

	params := mux.Vars(r)
	hash := params["hash"]

	q := query.NewGetSecretQuery(h.vault, hash)

	storedSecret, err := q.Execute()
	if err != nil {
		// TODO! Catch specificallt not found error
		http.Error(w, "", http.StatusNotFound)
		return
	}

	response := viewSecretResponseFromSecret(*storedSecret)

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
