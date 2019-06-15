package persistence

import (
	"errors"

	"github.com/vonhraban/secret-server/secret"
)

type InMemoryVault struct {
	storage map[string]*secret.Secret
}

func NewInMemoryVault() *InMemoryVault {
	return &InMemoryVault{
		storage: make(map[string]*secret.Secret),
	}
}

func (v *InMemoryVault) Store(secret *secret.Secret) (string, error) {
	// Errors?
	v.storage[secret.Hash] = secret

	return secret.Hash, nil
}

// TODO! Should I use UUID instead of string?
func (v *InMemoryVault) Retrieve(hash string) (*secret.Secret, error) {
	// TODO! Custom errors
	if val, ok := v.storage[hash]; ok {
		val.RemainingViews--
		return val, nil
	}

	return nil, errors.New("Not found")
}
