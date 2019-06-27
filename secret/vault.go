package secret

import (
	"github.com/pkg/errors"
)

var SecretNotFoundError = errors.New("Secret not found")
 
type Vault interface {
	Store(secret *Secret) error
	Retrieve(hash string) (*Secret, error)
	DecreaseRemainingViews(hash string) error
}
