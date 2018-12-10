package scheduler

import (
	"errors"
	"github.com/d1slike/go-sched/store"
)

type Option func(s *scheduler)

type Scheduler interface {
	Start()
	Shutdown() error
	RegisterExecutor(jType string, executor JobExecutor) Scheduler
	UnregisterExecutor(jType string)
	ScheduleJob(job MutableJob, trigger MutableTrigger) error
	GetJob(jKey string) (ImmutableJob, error)
	GetTrigger(tKey string) (ImmutableTrigger, error)
	DeleteJob(jKey string) (bool, error)
	DeleteTrigger(tKey string) (bool, error)
	GetJobs() ([]ImmutableJob, error)
	GetTriggers() ([]ImmutableTrigger, error)
}

type scheduler struct {
	name     string
	store    Store
	registry executorRegistry
}

func (s *scheduler) GetJob(jKey string) (ImmutableJob, error) {
	return s.store.GetJob(s.name, jKey)
}

func (s *scheduler) GetTrigger(tKey string) (ImmutableTrigger, error) {
	return s.store.GetTrigger(s.name, tKey)
}

func (s *scheduler) DeleteJob(jKey string) (bool, error) {
	return s.store.DeleteJob(s.name, jKey)
}

func (s *scheduler) DeleteTrigger(tKey string) (bool, error) {
	return s.store.DeleteTrigger(s.name, tKey)
}

func (s *scheduler) GetJobs() ([]ImmutableJob, error) {
	return s.store.GetJobs(s.name)
}

func (s *scheduler) GetTriggers() ([]ImmutableTrigger, error) {
	return s.store.GetTriggers(s.name)
}

func (s *scheduler) UnregisterExecutor(jType string) {
	s.registry.Unregister(jType)
}

func (s *scheduler) Start() {

}

func (s *scheduler) Shutdown() error {

}

func (s *scheduler) RegisterExecutor(jType string, executor JobExecutor) Scheduler {
	s.registry.Register(jType, executor)
	return s
}

func (s *scheduler) ScheduleJob(job MutableJob, tri MutableTrigger) error {
	j, err := job.ToImmutable()
	if err != nil {
		return err
	}

	t, err := tri.ToImmutable()
	if err != nil {
		return err
	}

	trigger, ok := t.(*trigger)
	if !ok {
		return errors.New("unknown trigger type")
	}
	trigger.jobKey = j.Key()
	trigger.state = StateScheduled

	if err := s.store.InsertJob(s.name, j); err != nil {
		return err
	}

	if err := s.store.InsertTrigger(s.name, trigger); err != nil {
		return err
	}

	return nil
}

func NewScheduler(name string, opts ...Option) Scheduler {
	s := &scheduler{
		name:     name,
		registry: newDefaultExecutorRegistry(),
	}

	for _, o := range opts {
		o(s)
	}

	if s.store == nil {
		s.store = store.NewInMemoryStore()
	}

	return s
}

func WithStore(store Store) Option {
	return func(s *scheduler) {
		s.store = store
	}
}

func WithExecutors(m map[string]JobExecutor) Option {
	return func(s *scheduler) {
		s.registry.RegisterAll(m)
	}
}
