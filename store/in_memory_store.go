package store

import "github.com/d1slike/go-sched"

type inMemoryStore struct {
}

func NewInMemoryStore() scheduler.Store {
	return &inMemoryStore{}
}
