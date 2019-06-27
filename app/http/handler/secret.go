package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/vonhraban/secret-server/core/log"
	"github.com/vonhraban/secret-server/secret"
	"github.com/vonhraban/secret-server/secret/cmd"
	"github.com/vonhraban/secret-server/secret/query"
)

type SecretHandler struct {
	vault  secret.Vault
	clock  secret.Clock
	logger log.Logger
}

func NewSecretHandler(
	vault secret.Vault,
	clock secret.Clock,
	logger log.Logger,
) *SecretHandler {
	return &SecretHandler{
		vault:  vault,
		clock:  clock,
		logger: logger,
	}
}

func (h *SecretHandler) Persist(w http.ResponseWriter, r *http.Request) {
	request, err := buildPersistSecretRequestFromHTTPRequest(r)
	if err != nil {
		var response interface{}
		switch err.(type) {
		case *EmptyValueError:
			h.logger.Warningf("Validation error %s", err)
			response = NewErrorResponse(err.Error())
			respond(w, r, h.logger, response)

		default:
			h.logger.Error(err)
			response = NewErrorResponse("Internal Error")
			respond(w, r, h.logger, response)
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
		respond(w, r, h.logger, response)

		return
	}

	q := query.NewGetSecretQuery(h.vault, hash)
	storedSecret, err := q.Execute()
	if err != nil {
		h.logger.Error(err)
		response := NewErrorResponse("Internal Error")
		respond(w, r, h.logger, response)

		return
	}

	response := buildPersistSecretResponseFromSecret(*storedSecret)
	respond(w, r, h.logger, response)
}

func (h *SecretHandler) View(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hash := params["hash"]

	q := query.NewGetSecretQuery(h.vault, hash)

	storedSecret, err := q.Execute()
	if err != nil {
		if err == secret.SecretNotFoundError {
			http.Error(w, "", http.StatusNotFound)
			return	
		}

		h.logger.Error(err)
		response := NewErrorResponse("Internal Error")
		respond(w, r, h.logger, response)
	}

	response := buildViewSecretResponseFromSecret(*storedSecret)

	decreaseViewsCmd := cmd.NewDecreaseRemainingViewsCommand(h.vault, hash)
	if err := decreaseViewsCmd.Execute(); err != nil {
		h.logger.Errorf("Could not decrease the number of avaialble views for a secret %s", hash)
	}

	respond(w, r, h.logger, response)
}
