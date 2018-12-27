package scheduler

import (
	"context"
	"github.com/d1slike/go-sched/internal"
	"github.com/d1slike/go-sched/jobs"
	"github.com/d1slike/go-sched/stores"
	"github.com/d1slike/go-sched/triggers"
)

type Option func(s *scheduler)

type Scheduler interface {
	Start()
	Shutdown(ctx context.Context) error
	RegisterExecutor(jType string, executor JobExecutor) Scheduler
	UnregisterExecutor(jType string)
	ScheduleJob(job jobs.MutableJob, trigger triggers.MutableTrigger) error
	GetJob(jKey string) (jobs.ImmutableJob, error)
	GetTrigger(tKey string) (triggers.ImmutableTrigger, error)
	DeleteJob(jKey string) (bool, error)
	DeleteTrigger(tKey string) (bool, error)
	GetJobs() ([]jobs.ImmutableJob, error)
	GetTriggers() ([]triggers.ImmutableTrigger, error)
	UpdateJob(job jobs.MutableJob) error
}

type scheduler struct {
	name     string
	store    stores.Store
	registry executorRegistry
	executor executor
	timers   Timers
}

func (s *scheduler) GetJob(jKey string) (jobs.ImmutableJob, error) {
	return s.store.GetJob(s.name, jKey)
}

func (s *scheduler) GetTrigger(tKey string) (triggers.ImmutableTrigger, error) {
	return s.store.GetTrigger(s.name, tKey)
}

func (s *scheduler) DeleteJob(jKey string) (bool, error) {
	ok, err := s.store.DeleteJob(s.name, jKey)
	if err != nil {
		return false, err
	}
	if ok {
		keys, err := s.store.DeleteTriggersByJobKey(s.name, jKey)
		if err != nil {
			return false, err
		}
		s.executor.CancelTriggers(keys...)
	}

	return ok, nil
}

func (s *scheduler) DeleteTrigger(tKey string) (bool, error) {
	ok, err := s.store.DeleteTrigger(s.name, tKey)
	if err != nil {
		return false, err
	}
	if ok {
		s.executor.CancelTriggers(tKey)
	}
	return ok, nil
}

func (s *scheduler) GetJobs() ([]jobs.ImmutableJob, error) {
	return s.store.GetJobs(s.name)
}

func (s *scheduler) GetTriggers() ([]triggers.ImmutableTrigger, error) {
	return s.store.GetTriggers(s.name)
}

func (s *scheduler) UnregisterExecutor(jType string) {
	s.registry.Unregister(jType)
}

func (s *scheduler) Start() {
	s.executor.Start()
}

func (s *scheduler) Shutdown(ctx context.Context) error {
	return s.executor.Shutdown(ctx)
}

func (s *scheduler) RegisterExecutor(jType string, executor JobExecutor) Scheduler {
	s.registry.Register(jType, executor)
	return s
}

func (s *scheduler) ScheduleJob(job jobs.MutableJob, tri triggers.MutableTrigger) error {
	j, err := job.ToImmutable()
	if err != nil {
		return err
	}

	t, err := tri.ToImmutable()
	if err != nil {
		return err
	}

	t = internal.ModifyTrigger(t, func(tr *internal.Trigger) {
		tr.TjobKey = j.Key()
		tr.Tstate = triggers.StateScheduled
	})

	if err := s.store.InsertJob(s.name, j); err != nil {
		return err
	}

	if err := s.store.InsertTrigger(s.name, t); err != nil {
		return err
	}

	return nil
}

func (s *scheduler) UpdateJob(job jobs.MutableJob) error {
	j, err := job.ToImmutable()
	if err != nil {
		return err
	}

	return s.store.UpdateJob(s.name, j)
}

func NewScheduler(name string, opts ...Option) Scheduler {
	s := &scheduler{
		name:     name,
		registry: newDefaultExecutorRegistry(),
		timers:   NewDefaultTimers(),
	}

	for _, o := range opts {
		o(s)
	}

	if s.store == nil {
		s.store = stores.NewInMemoryStore()
	}

	s.executor = newDefaultRuntimeExecutor(
		s.name,
		s.store,
		s.registry,
		s.timers,
	)

	return s
}

func WithStore(store stores.Store) Option {
	return func(s *scheduler) {
		s.store = store
	}
}

func WithExecutors(m map[string]JobExecutor) Option {
	return func(s *scheduler) {
		s.registry.RegisterAll(m)
	}
}

func WithTimers(timers Timers) Option {
	return func(s *scheduler) {
		s.timers = SetDefault(timers)
	}
}
