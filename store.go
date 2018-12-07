package scheduler

type Store interface {
	InsertJob(job ImmutableJob) error
	InsertTrigger(trigger ImmutableTrigger) error
	GetJob(jKey string) (ImmutableJob, bool)
	GetTrigger(tKey string) (ImmutableTrigger, bool)
	DeleteJob(jKey string) bool
	DeleteTrigger(tKey string) bool
	GetJobs() []ImmutableJob
	GetTriggers() []ImmutableTrigger
}
