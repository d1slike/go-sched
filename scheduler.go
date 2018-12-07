package scheduler

import (
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
	store    store.Store
	registry executorRegistry
}

func (s *scheduler) Start() {

}

func (s *scheduler) Shutdown() error {

}

func (s *scheduler) RegisterExecutor(jType string, executor JobExecutor) {
	s.registry.Register(jType, executor)
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

func WithStore(store store.Store) Option {
	return func(s *scheduler) {
		s.store = store
	}
}
