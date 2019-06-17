package persistence

import (
	"errors"

	"github.com/vonhraban/secret-server/secret"
)

type InMemoryVault struct {
	storage map[string]*secret.Secret
	clock   secret.Clock
}

func NewInMemoryVault(clock secret.Clock) *InMemoryVault {
	return &InMemoryVault{
		storage: make(map[string]*secret.Secret),
		clock:   clock,
	}
}

func (v *InMemoryVault) Store(secret *secret.Secret) (string, error) {
	// Errors?
	v.storage[secret.Hash] = secret

	return secret.Hash, nil
}

func (v *InMemoryVault) Retrieve(hash string) (*secret.Secret, error) {
	// TODO! Custom errors
	if val, ok := v.storage[hash]; ok && val.CanBeSeen(v.clock.GetCurrentTime()) {
		return val, nil
	}

	return nil, errors.New("Not found")
}

func (v *InMemoryVault) DecreaseRemainingViews(hash string) error {
	if val, ok := v.storage[hash]; ok {
		val.RemainingViews--
		return nil
	}

	return errors.New("Not found")
}
