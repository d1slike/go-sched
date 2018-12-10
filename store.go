package scheduler

type Store interface {
	InsertJob(sName string, job ImmutableJob) error
	InsertTrigger(sName string, trigger ImmutableTrigger) error
	GetJob(sName string, jKey string) (ImmutableJob, error)
	GetTrigger(sName string, tKey string) (ImmutableTrigger, error)
	DeleteJob(sName string, jKey string) (bool, error)
	DeleteTrigger(sName string, tKey string) (bool, error)
	GetJobs(sName string) ([]ImmutableJob, error)
	GetTriggers(sName string) ([]ImmutableTrigger, error)
	AcquireTriggers(sName string) ([]ImmutableTrigger, error)
	UpdateTrigger(sName string, trigger ImmutableTrigger) error
	UpdateJob(sName string, job ImmutableJob) error
	DeleteExhaustedTriggers(sName string) (int, error)
}
