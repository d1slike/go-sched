package stores

import (
	"github.com/d1slike/go-sched/internal"
	"github.com/d1slike/go-sched/jobs"
	"github.com/d1slike/go-sched/triggers"
	"strings"
	"sync"
)

type inMemoryStore struct {
	tLock sync.RWMutex
	jLock sync.RWMutex

	tMap map[string]triggers.ImmutableTrigger
	jMap map[string]jobs.ImmutableJob
}

func (s *inMemoryStore) InsertJob(sName string, job jobs.ImmutableJob) error {
	s.jLock.Lock()
	defer s.jLock.Unlock()

	if _, exists := s.jMap[key(sName, job.Key())]; exists {
		return ErrJobAlreadyExists
	}

	s.jMap[key(sName, job.Key())] = job

	return nil
}

func (s *inMemoryStore) InsertTrigger(sName string, trigger triggers.ImmutableTrigger) error {
	s.tLock.Lock()
	defer s.tLock.Unlock()

	if _, exists := s.tMap[key(sName, trigger.Key())]; exists {
		return ErrTriggerAlreadyExists
	}

	s.tMap[key(sName, trigger.Key())] = trigger

	return nil
}

func (s *inMemoryStore) GetJob(sName string, jKey string) (jobs.ImmutableJob, error) {
	s.jLock.RLock()
	defer s.jLock.RUnlock()

	j, ok := s.jMap[key(sName, jKey)]
	if !ok {
		return nil, nil
	}

	return j, nil
}

func (s *inMemoryStore) GetTrigger(sName string, tKey string) (triggers.ImmutableTrigger, error) {
	s.tLock.RLock()
	defer s.tLock.RUnlock()

	t, ok := s.tMap[key(sName, tKey)]
	if !ok {
		return nil, nil
	}

	return t, nil
}

func (s *inMemoryStore) DeleteJob(sName string, jKey string) (bool, error) {
	s.jLock.Lock()
	defer s.jLock.Unlock()

	_, ok := s.jMap[key(sName, jKey)]
	delete(s.jMap, key(sName, jKey))

	return ok, nil
}

func (s *inMemoryStore) DeleteTrigger(sName string, tKey string) (bool, error) {
	s.tLock.Lock()
	defer s.tLock.Unlock()

	_, ok := s.tMap[key(sName, tKey)]
	delete(s.tMap, key(sName, tKey))

	return ok, nil
}

func (s *inMemoryStore) GetJobs(sName string) ([]jobs.ImmutableJob, error) {
	s.jLock.RLock()
	defer s.jLock.RUnlock()

	arr := make([]jobs.ImmutableJob, 0, len(s.jMap))
	for key, job := range s.jMap {
		if strings.HasPrefix(key, sName) {
			arr = append(arr, job)
		}
	}

	return arr, nil
}

func (s *inMemoryStore) GetTriggers(sName string) ([]triggers.ImmutableTrigger, error) {
	s.tLock.RLock()
	defer s.tLock.RUnlock()

	arr := make([]triggers.ImmutableTrigger, 0, len(s.jMap))
	for key, trigger := range s.tMap {
		if strings.HasPrefix(key, sName) {
			arr = append(arr, trigger)
		}
	}

	return arr, nil
}

func (s *inMemoryStore) AcquireTriggers(sName string) ([]triggers.ImmutableTrigger, error) {
	s.tLock.Lock()
	defer s.tLock.Unlock()

	arr := make([]triggers.ImmutableTrigger, 0, len(s.jMap)/2)
	for key, trigger := range s.tMap {
		if strings.HasPrefix(key, sName) && trigger.State() == triggers.StateScheduled {
			trigger = internal.ModifyTrigger(trigger, func(tr *internal.Trigger) {
				tr.TState = triggers.StateAcquired
			})
			s.tMap[key] = trigger
			arr = append(arr, trigger)
		}
	}

	return arr, nil
}

func (s *inMemoryStore) UpdateTrigger(sName string, trigger triggers.ImmutableTrigger) error {
	s.tLock.Lock()
	defer s.tLock.Unlock()

	_, ok := s.tMap[key(sName, trigger.Key())]
	if !ok {
		return ErrTriggerNotFound
	}

	s.tMap[key(sName, trigger.Key())] = trigger

	return nil
}

func (s *inMemoryStore) UpdateJob(sName string, job jobs.ImmutableJob) error {
	s.jLock.Lock()
	defer s.jLock.Unlock()

	_, ok := s.jMap[key(sName, job.Key())]
	if !ok {
		return ErrJobNotFound
	}

	s.jMap[key(sName, job.Key())] = job

	return nil
}

func (s *inMemoryStore) DeleteExhaustedTriggers(sName string) (int, error) {
	s.tLock.Lock()
	defer s.tLock.Unlock()

	deleted := 0
	newMap := make(map[string]triggers.ImmutableTrigger, len(s.tMap))
	for key, trigger := range s.tMap {
		isOwner := strings.HasPrefix(key, sName)
		if (isOwner && trigger.State() != triggers.StateExhausted) || !isOwner {
			newMap[key] = trigger
		} else if isOwner {
			deleted++
		}
	}

	s.tMap = newMap

	return deleted, nil
}

func NewInMemoryStore() Store {
	return &inMemoryStore{
		tMap: make(map[string]triggers.ImmutableTrigger),
		jMap: make(map[string]jobs.ImmutableJob),
	}
}

func key(sName, entityKey string) string {
	return sName + entityKey
}
