package scheduler

import (
	"github.com/d1slike/go-sched/internal"
	"github.com/d1slike/go-sched/jobs"
)

func NewJob() jobs.MutableJob {
	return internal.NewJob()
}
