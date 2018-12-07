package scheduler

type JobExecutor func(ctx JobContext) error

type executorRegistry interface {
	Register(jType string, executor JobExecutor)
	RegisterAll(map[string]JobExecutor)
	Unregister(jType string)
	GetExecutor(jType string) (JobExecutor, bool)
}

type defaultExecutorRegistry struct {
	registry map[string]JobExecutor
}

func (r *defaultExecutorRegistry) Register(jType string, executor JobExecutor) {
	r.registry[jType] = executor
}

func (r *defaultExecutorRegistry) RegisterAll(m map[string]JobExecutor) {
	for t, e := range m {
		r.registry[t] = e
	}
}

func (r *defaultExecutorRegistry) Unregister(jType string) {
	delete(r.registry, jType)
}

func (r *defaultExecutorRegistry) GetExecutor(jType string) (JobExecutor, bool) {
	e, ok := r.registry[jType]
	return e, ok
}

func newDefaultExecutorRegistry() executorRegistry {
	return &defaultExecutorRegistry{
		registry: make(map[string]JobExecutor),
	}
}
