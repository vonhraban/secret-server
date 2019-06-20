package query

import (
	"github.com/pkg/errors"
	"github.com/vonhraban/secret-server/secret"
)

type GetSecretQuery struct {
	vault secret.Vault
	hash  string
}

func NewGetSecretQuery(vault secret.Vault, hash string) *GetSecretQuery {
	return &GetSecretQuery{
		vault: vault,
		hash:  hash,
	}
}

func (q *GetSecretQuery) Execute() (*secret.Secret, error) {
	value, err := q.vault.Retrieve(q.hash)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not query a secret %s", q.hash)

	}

	return value, nil
}
