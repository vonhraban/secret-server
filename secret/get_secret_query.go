package secret

type GetSecretQuery struct{}

func (q *GetSecretQuery) Execute(vault Vault, hash string) (*Secret, error){
	value, err := vault.Retrieve(hash)
	if err != nil {
		// TODO Wrapf
		return nil, err
	}

	return value, nil
}