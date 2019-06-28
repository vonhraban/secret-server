package persistence

import (
	"errors"
	"sync"

	"github.com/vonhraban/secret-server/secret"
)

type inMemoryVault struct {
	storage map[string]*secret.Secret
	mux     sync.Mutex
}

func NewInMemoryVault() *inMemoryVault {
	return &inMemoryVault{
		storage: make(map[string]*secret.Secret),
	}
}

func (v *inMemoryVault) Store(secret *secret.Secret) error {
	v.storage[secret.Hash] = secret

	return nil
}

func (v *inMemoryVault) Retrieve(hash string) (*secret.Secret, error) {
	if val, ok := v.storage[hash]; ok {
		return val, nil
	}

	return nil, secret.SecretNotFoundError
}

func (v *inMemoryVault) DecreaseRemainingViews(hash string) error {
	if val, ok := v.storage[hash]; ok {
		v.mux.Lock()
		val.RemainingViews--
		v.mux.Unlock()
		return nil
	}

	return errors.New("Not found")
}
