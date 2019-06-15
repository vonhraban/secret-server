package secret

import "time"

type AddSecret struct{}

func (cmd *AddSecret) Execute(vault Vault, clock Clock, secretText string, maxViews int) (string, error) {
	// TODO! Validate max views is greater than 0
	now := clock.GetCurrentTime()
	//panic(now)
	expirationTime := now.Add(time.Hour * 1) // 1 hour expiration
	secret := &Secret{
		SecretText:     secretText,
		RemainingViews: maxViews,
		CreatedAt:      clock.GetCurrentTime(),
		ExpiresAt:      expirationTime,
	}

	hash, err := vault.Store(secret)
	if err != nil {
		// todo! errwrapf
		return "", err
	}

	return hash, nil
}
