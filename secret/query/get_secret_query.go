package query

import (
	"github.com/pkg/errors"
	"github.com/vonhraban/secret-server/secret"
)

type getSecretQuery struct {
	vault secret.Vault
	hash  string
}

func NewGetSecretQuery(vault secret.Vault, hash string) *getSecretQuery {
	return &getSecretQuery{
		vault: vault,
		hash:  hash,
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

	return value, nil
}
