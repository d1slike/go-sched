package scheduler

import (
	"github.com/d1slike/go-sched/triggers"
	"github.com/d1slike/go-sched/utils"
	"time"
)

type future struct {
	timer    *time.Timer
	t        triggers.ImmutableTrigger
	running  *utils.AtomicBool
	canceled *utils.AtomicBool
}

func (f *future) Cancel() {
	f.canceled.Set(true)
	f.timer.Stop()
}

func (f *future) IsRunning() bool {
	return f.running.Get()
}

func (f *future) IsCanceled() bool {
	return f.canceled.Get()
}

func (f *future) Run() {
	f.running.Set(true)
}
