package query

import "github.com/vonhraban/secret-server/secret"

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
		// TODO Wrapf
		return nil, err
	}

	return value, nil
}
