package secret

type Vault interface {
	Store(secret *Secret) (string, error)
	Retrieve(hash string) (*Secret, error)
}
