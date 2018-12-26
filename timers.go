package scheduler

import "time"

const (
	DefaultTriggerStealTimeout = 1 * time.Second
)

type Timers struct {
	TriggerStealTimeout time.Duration
}

func NewDefaultTimers() Timers {
	return SetDefault(Timers{})
}

func SetDefault(t Timers) Timers {
	if t.TriggerStealTimeout <= 0 {
		t.TriggerStealTimeout = DefaultTriggerStealTimeout
	}

	return t
}
