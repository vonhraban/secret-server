package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/vonhraban/secret-server/secret"
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

	cmd := secret.NewAddSecretCommand(
		h.vault,
		h.clock,
		hash,
		request.secret,
		request.expireAfterViews,
		request.expireAfter,
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	query := secret.NewGetSecretQuery(h.vault, hash)
	storedSecret, err := query.Execute()
	if err != nil {
		panic(err)
	}

	response := persistSecretResponseFromSecret(*storedSecret)

	json.NewEncoder(w).Encode(response)
}

func (h *secretHandler) View(w http.ResponseWriter, r *http.Request) {
	// TODO! Validation

	params := mux.Vars(r)
	hash := params["hash"]

	cmd := secret.NewGetSecretQuery(h.vault, hash)

	storedSecret, err := cmd.Execute()
	if err != nil {
		// TODO! Catch specificallt not found error
		http.Error(w, "", http.StatusNotFound)
		return
	}

	response := viewSecretResponseFromSecret(*storedSecret)

	decreaseViewsCmd := secret.NewDecreaseRemainingViewsCommand(h.vault, hash)
	if err := decreaseViewsCmd.Execute(); err != nil {
		panic(err)
	}

	// Since we now decreased the number of available views in a store secret, we need to descrease it in the response too
	// and we want to do it only if there are more that 0 remaining views
	if response.RemainingViews > 0 {
		response.RemainingViews--
	}

	json.NewEncoder(w).Encode(response)
}
