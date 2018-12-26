package triggers

import (
	"errors"
	"time"
)

const (
	RepeatInfinity = Repeats(-1)
)

const (
	StateScheduled = TriggerState("SCHEDULED")
	StateAcquired  = TriggerState("ACQUIRED")
	StateExhausted = TriggerState("EXHAUSTED")
)

var (
	ErrEmptyTriggerKey  = errors.New("empty trigger key")
	ErrEmptyCronSpec    = errors.New("empty cron specification")
	ErrAlreadyExhausted = errors.New("trigger already exhausted")
	ErrInvalidLocation  = "invalid location: %v"
	ErrInvalidCronSpec  = "invalid cron spec: %v"
)

type Repeats int

type TriggerState string

type MutableTrigger interface {
	WithKey(tKey string) MutableTrigger
	WithFromTime(from time.Time) MutableTrigger
	WithToTime(to time.Time) MutableTrigger
	WithRepeats(repeats Repeats) MutableTrigger
	WithCron(spec string) MutableTrigger
	InLocation(loc string) MutableTrigger
	ToImmutable() (ImmutableTrigger, error)
}

type ImmutableTrigger interface {
	Key() string
	JobKey() string
	FromTime() *time.Time
	ToTime() *time.Time
	Repeats() Repeats
	CronSpec() string
	Location() *time.Location
	State() TriggerState
	NextTriggerTime() time.Time
}

func Repeat(count int) Repeats {
	return Repeats(count)
}
