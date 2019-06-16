package secret

type GetSecretQuery struct {
	vault Vault
	hash  string
}

func NewGetSecretQuery(vault Vault, hash string) *GetSecretQuery {
	return &GetSecretQuery{
		vault: vault,
		hash:  hash,
	}
}

func (q *GetSecretQuery) Execute() (*Secret, error) {
	value, err := q.vault.Retrieve(q.hash)
	if err != nil {
		// TODO Wrapf
		return nil, err
	}

	return value, nil
}
