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
			respond(w, r, response)

		default:
			h.logger.Error(err)
			response = NewErrorResponse("Internal Error")
			respond(w, r, response)
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
		respond(w, r, response)

		return
	}

	q := query.NewGetSecretQuery(h.vault, hash)
	storedSecret, err := q.Execute()
	if err != nil {
		h.logger.Error(err)
		response := NewErrorResponse("Internal Error")
		respond(w, r, response)

		return
	}

	response := buildPersistSecretResponseFromSecret(*storedSecret)
	respond(w, r, response)
}

func (h *SecretHandler) View(w http.ResponseWriter, r *http.Request) {
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

	respond(w, r, response)
}
