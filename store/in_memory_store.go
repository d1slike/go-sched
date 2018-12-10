package store

import (
	"github.com/d1slike/go-sched"
	"sync"
)

type inMemoryStore struct {
	tLock sync.RWMutex
	jLock sync.RWMutex
}

func (*inMemoryStore) InsertJob(sName string, job scheduler.ImmutableJob) error {
	panic("implement me")
}

func (*inMemoryStore) InsertTrigger(sName string, trigger scheduler.ImmutableTrigger) error {
	panic("implement me")
}

func (*inMemoryStore) GetJob(sName string, jKey string) (scheduler.ImmutableJob, error) {
	panic("implement me")
}

func (*inMemoryStore) GetTrigger(sName string, tKey string) (scheduler.ImmutableTrigger, error) {
	panic("implement me")
}

func (*inMemoryStore) DeleteJob(sName string, jKey string) (bool, error) {
	panic("implement me")
}

func (*inMemoryStore) DeleteTrigger(sName string, tKey string) (bool, error) {
	panic("implement me")
}

func (*inMemoryStore) GetJobs(sName string) ([]scheduler.ImmutableJob, error) {
	panic("implement me")
}

func (*inMemoryStore) GetTriggers(sName string) ([]scheduler.ImmutableTrigger, error) {
	panic("implement me")
}

func (*inMemoryStore) AcquireTriggers(sName string) ([]scheduler.ImmutableTrigger, error) {
	panic("implement me")
}

func (*inMemoryStore) UpdateTrigger(sName string, trigger scheduler.ImmutableTrigger) error {
	panic("implement me")
}

func (*inMemoryStore) UpdateJob(sName string, job scheduler.ImmutableJob) error {
	panic("implement me")
}

func (*inMemoryStore) DeleteExhaustedTriggers(sName string) (int, error) {
	panic("implement me")
}

func NewInMemoryStore() scheduler.Store {
	return &inMemoryStore{}
}
