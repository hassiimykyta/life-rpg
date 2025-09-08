package password

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}
