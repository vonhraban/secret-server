package secret

type Vault interface {
	Store(secret *Secret) error
	Retrieve(hash string) (*Secret, error)
	DecreaseRemainingViews(hash string) error
}
