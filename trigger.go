package scheduler

import (
	"github.com/d1slike/go-sched/internal"
	"github.com/d1slike/go-sched/triggers"
)

func NewTrigger() triggers.MutableTrigger {
	return internal.NewTrigger()
}
