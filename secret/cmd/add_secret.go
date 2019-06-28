package cmd

import (
	"time"

	"github.com/pkg/errors"
	"github.com/vonhraban/secret-server/secret"
)

type addSecret struct {
	vault      secret.Vault
	clock      secret.Clock
	hash       string
	secretText string
	maxViews   int
	ttlMins    int
}

func NewAddSecretCommand(vault secret.Vault, clock secret.Clock, hash string, secretText string, maxViews int, ttlMins int) *addSecret {
	return &addSecret{
		vault:      vault,
		clock:      clock,
		hash:       hash,
		secretText: secretText,
		maxViews:   maxViews,
		ttlMins:    ttlMins,
	}
}

func (cmd *addSecret) Execute() error {
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

	if err := cmd.vault.Store(secret); err != nil {
		return errors.Wrapf(err, "Could not store secret %s", secret.Hash)
	}

	return nil
}
