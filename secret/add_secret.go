package secret

import "time"

type AddSecret struct{}

func (cmd *AddSecret) Execute(vault Vault, token string) (string, error) {
	secret := &Secret{
		Token:     token,
		MaxUses:   5,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
	}

	id, err := vault.Store(secret)
	if err != nil {
		// todo! errwrapf
		return "", err
	}

	return id, nil
}
