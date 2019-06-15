package secret

import "time"

type AddSecret struct{}

func (cmd *AddSecret) Execute(vault Vault, secretText string, maxViews int) (string, error) {
	// TODO! Validate max views is greater than 0
	secret := &Secret{
		SecretText:     secretText,
		RemainingViews: maxViews,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now(),
	}

	hash, err := vault.Store(secret)
	if err != nil {
		// todo! errwrapf
		return "", err
	}

	return hash, nil
}
