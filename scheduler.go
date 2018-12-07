package scheduler

import (
	"errors"
	"github.com/d1slike/go-sched/store"
)

type Option func(s *scheduler)

type Scheduler interface {
	Start()
	Shutdown() error
	RegisterExecutor(jType string, executor JobExecutor)
	ScheduleJob(job MutableJob, trigger MutableTrigger)
}

type scheduler struct {
	name     string
	store    Store
	registry executorRegistry
}

func (s *scheduler) Start() {

}

func (s *scheduler) Shutdown() error {

}

func (s *scheduler) RegisterExecutor(jType string, executor JobExecutor) {
	s.registry.Register(jType, executor)
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
