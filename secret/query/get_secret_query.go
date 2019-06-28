package query

import (
	"github.com/pkg/errors"
	"github.com/vonhraban/secret-server/secret"
)

type getSecretQuery struct {
	vault secret.Vault
	hash  string
	clock secret.Clock
}

func NewGetSecretQuery(vault secret.Vault, hash string, clock secret.Clock) *getSecretQuery {
	return &getSecretQuery{
		vault: vault,
		hash:  hash,
		clock: clock,
	}
}

func (q *getSecretQuery) Execute() (*secret.Secret, error) {
	value, err := q.vault.Retrieve(q.hash)

	if err != nil {
		if err == secret.SecretNotFoundError {
			return nil, err
		}

		return nil, errors.Wrapf(err, "An error occured querying a secret %s", q.hash)
	}

	// The value has been retrieved, but it is accessuble?
	if !value.CanBeSeen(q.clock.GetCurrentTime()) {
		return nil, secret.SecretNotFoundError
	}

	return value, nil
}
