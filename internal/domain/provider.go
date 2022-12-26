package domain

// IDProvider represent id generator contract
type IDProvider interface {
	Generate() string
}

// HashProvider represent hasher contract
type HashProvider interface {
	Hash(raw string) ([]byte, error)
}
