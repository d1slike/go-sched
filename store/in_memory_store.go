package store

type inMemoryStore struct {
}

func NewInMemoryStore() Store {
	return &inMemoryStore{}
}
