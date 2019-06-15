package secret

type Vault interface {
	Store(secret *Secret) (string, error)
	Retrieve(id string) (*Secret, error)
}
