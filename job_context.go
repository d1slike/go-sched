package scheduler

import (
	"errors"
	"github.com/d1slike/go-sched/jobs"
	"github.com/d1slike/go-sched/json"
	"github.com/d1slike/go-sched/triggers"
)

var (
	ErrNoData = errors.New("no data")
)

type JobContext interface {
	Trigger() triggers.ImmutableTrigger
	Job() jobs.ImmutableJob
	UnmarshalJobData(ptr interface{}) error
	UnmarshalTriggerData(ptr interface{}) error
}

type jobCtx struct {
	job     jobs.ImmutableJob
	trigger triggers.ImmutableTrigger
}

func (ctx *jobCtx) Trigger() triggers.ImmutableTrigger {
	return ctx.trigger
}

func (ctx *jobCtx) Job() jobs.ImmutableJob {
	return ctx.job
}

func (ctx *jobCtx) UnmarshalJobData(ptr interface{}) error {
	if len(ctx.job.Data()) == 0 {
		return ErrNoData
	}
	return json.Provider.Unmarshal(ctx.job.Data(), ptr)
}

func (ctx *jobCtx) UnmarshalTriggerData(ptr interface{}) error {
	if len(ctx.trigger.Data()) == 0 {
		return ErrNoData
	}
	return json.Provider.Unmarshal(ctx.trigger.Data(), ptr)
}
