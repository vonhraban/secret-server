package persistence

import (
	"errors"

	"github.com/vonhraban/secret-server/secret"
)

type inMemoryVault struct {
	storage map[string]*secret.Secret
	clock   secret.Clock
}

func NewInMemoryVault(clock secret.Clock) *inMemoryVault {
	return &inMemoryVault{
		storage: make(map[string]*secret.Secret),
		clock:   clock,
	}
}

func (v *inMemoryVault) Store(secret *secret.Secret) error {
	v.storage[secret.Hash] = secret

	return nil
}

func (v *inMemoryVault) Retrieve(hash string) (*secret.Secret, error) {
	// TODO! Custom errors
	if val, ok := v.storage[hash]; ok && val.CanBeSeen(v.clock.GetCurrentTime()) {
		return val, nil
	}

	return nil, secret.SecretNotFoundError
}

func (v *inMemoryVault) DecreaseRemainingViews(hash string) error {
	if val, ok := v.storage[hash]; ok {
		val.RemainingViews--
		return nil
	}

	return errors.New("Not found")
}
