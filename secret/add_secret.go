package secret

import "time"

type AddSecret struct{}

func (cmd *AddSecret) Execute(vault Vault, clock Clock, secretText string, maxViews int, ttlMins int) (string, error) {
	// TODO! Validate max views is greater than 0
	now := clock.GetCurrentTime()
	var expirationTime time.Time
	if ttlMins != 0 {
		expirationTime = now.Add(time.Minute * time.Duration(ttlMins))
	}
	secret := &Secret{
		SecretText:     secretText,
		RemainingViews: maxViews,
		CreatedAt:      now,
		ExpiresAt:      expirationTime,
	}

	hash, err := vault.Store(secret)
	if err != nil {
		// todo! errwrapf
		return "", err
	}

	return hash, nil
}
