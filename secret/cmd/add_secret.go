package cmd

import (
	"time"
	"github.com/vonhraban/secret-server/secret"
)

type AddSecret struct{
	vault secret.Vault
	clock secret.Clock
	hash string
	secretText string
	maxViews int
	ttlMins int
}

func NewAddSecretCommand(vault secret.Vault, clock secret.Clock, hash string, secretText string, maxViews int, ttlMins int) *AddSecret {
	return &AddSecret{
		vault: vault,
		clock: clock,
		hash: hash,
		secretText: secretText,
		maxViews: maxViews,
		ttlMins: ttlMins,
	}
}

func (cmd *AddSecret) Execute() error {
	// TODO! Validate max views is greater than 0
	now := cmd.clock.GetCurrentTime()
	var expirationTime time.Time
	if cmd.ttlMins != 0 {
		expirationTime = now.Add(time.Minute * time.Duration(cmd.ttlMins))
	}
	secret := &secret.Secret{
		Hash:           cmd.hash,
		SecretText:     cmd.secretText,
		RemainingViews: cmd.maxViews,
		CreatedAt:      now,
		ExpiresAt:      expirationTime,
	}

	// TODO Do I even need the vault to return the ID?
	_, err := cmd.vault.Store(secret)
	if err != nil {
		// todo! errwrapf
		return  err
	}

	return nil
}
