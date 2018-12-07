package scheduler

type Store interface {
	SaveJob(job ImmutableJob) error
	SaveTrigger(trigger ImmutableTrigger) error
}
