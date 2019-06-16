package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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

type persistSecretRequest struct {
	secret           string
	expireAfterViews int
	expireAfter      int
}

func persistSecretRequestFromHTTPRequest(r *http.Request) (*persistSecretRequest, error) {
	r.ParseForm()
	secret := r.FormValue("secret")
	expireAfterViews, err := strconv.Atoi(r.FormValue("expireAfterViews"))
	if err != nil {
		return nil, err
	}

	expireAfter, err := strconv.Atoi(r.FormValue("expireAfter"))
	if err != nil {
		return nil, err
	}

	return &persistSecretRequest{
		secret:           secret,
		expireAfterViews: expireAfterViews,
		expireAfter:      expireAfter,
	}, nil
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

	json.NewEncoder(w).Encode(storedSecret)
}
