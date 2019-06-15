package persistence

import (
	"errors"

	uuid "github.com/satori/go.uuid"
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
	id := uuid.NewV4()
	// Errors?
	v.storage[id.String()] = secret

	return id.String(), nil
}

func (v *InMemoryVault) Retrieve(id string) (*secret.Secret, error) {
	// TODO! Increase the uses
	// TODO! Custom errors
	if val, ok := v.storage[id]; ok {
		return val, nil
	}

	return nil, errors.New("Not found")
}
