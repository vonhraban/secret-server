package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/vonhraban/secret-server/secret"
)

type SecretHandler struct {
	Vault secret.Vault
	Clock secret.Clock
}

type PersistSecretRequest struct {
	Secret           string `json:"`
	ExpireAfterViews int
	ExpireAfter      int
}

func PersistSecretRequestFromHttpRequest(r *http.Request) (*PersistSecretRequest, error) {
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

	return &PersistSecretRequest{
		Secret:           secret,
		ExpireAfterViews: expireAfterViews,
		ExpireAfter:      expireAfter,
	}, nil
}

func (h *SecretHandler) Persist(w http.ResponseWriter, r *http.Request) {
	//panic(fmt.Sprintf("%+v", r))
	request, err := PersistSecretRequestFromHttpRequest(r)
	if err != nil {
		panic(err)
	}
	// TODO! Validation

	hash := uuid.NewV4().String()

	cmd := secret.NewAddSecretCommand(
		h.Vault,
		h.Clock,
		hash,
		request.Secret,
		request.ExpireAfterViews,
		request.ExpireAfter,
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	query := secret.NewGetSecretQuery(h.Vault, hash)
	storedSecret, err := query.Execute()
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(storedSecret)
}
