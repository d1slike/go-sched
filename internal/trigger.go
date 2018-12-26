package internal

import (
	"fmt"
	"github.com/d1slike/go-sched/triggers"
	"github.com/robfig/cron"
	"time"
)

type Trigger struct {
	TKey      string
	TJobKey   string
	TFromTime *time.Time
	TToTime   *time.Time
	TRepeats  triggers.Repeats
	TCronSpec string
	TLocation string

	TState         triggers.TriggerState
	TLoc           *time.Location
	TTriggeredTime triggers.Repeats
	TSched         cron.Schedule
	TNextTime      time.Time
}

func (t *Trigger) State() triggers.TriggerState {
	return t.TState
}

func (t *Trigger) JobKey() string {
	return t.TKey
}

func (t *Trigger) WithKey(tKey string) triggers.MutableTrigger {
	t.TKey = tKey
	return t
}

func (t *Trigger) WithFromTime(from time.Time) triggers.MutableTrigger {
	t.TFromTime = &from
	return t
}

func (t *Trigger) WithToTime(to time.Time) triggers.MutableTrigger {
	t.TToTime = &to
	return t
}

func (t *Trigger) WithRepeats(repeats triggers.Repeats) triggers.MutableTrigger {
	t.TRepeats = repeats
	return t
}

func (t *Trigger) WithCron(spec string) triggers.MutableTrigger {
	t.TCronSpec = spec
	return t
}

func (t *Trigger) InLocation(loc string) triggers.MutableTrigger {
	t.TLocation = loc
	return t
}

func (t *Trigger) ToImmutable() (triggers.ImmutableTrigger, error) {
	if t.TKey == "" {
		return nil, triggers.ErrEmptyTriggerKey
	}
	if t.TCronSpec == "" {
		return nil, triggers.ErrEmptyCronSpec
	}

	loc, err := time.LoadLocation(t.TLocation)
	if err != nil {
		return nil, fmt.Errorf(triggers.ErrInvalidLocation, err)
	}
	t.TLoc = loc

	sched, err := cron.Parse(t.TCronSpec)
	if err != nil {
		return nil, fmt.Errorf(triggers.ErrInvalidCronSpec, err)
	}
	t.TSched = sched

	t.TNextTime = t.TSched.Next(time.Now().In(t.TLoc))

	if t.TNextTime.IsZero() {
		return nil, triggers.ErrAlreadyExhausted
	}

	return t, nil
}

func (t *Trigger) Key() string {
	return t.TKey
}

func (t *Trigger) FromTime() *time.Time {
	return t.TFromTime
}

func (t *Trigger) ToTime() *time.Time {
	return t.TToTime
}

func (t *Trigger) Repeats() triggers.Repeats {
	return t.TRepeats
}

func (t *Trigger) CronSpec() string {
	return t.TCronSpec
}

func (t *Trigger) Location() *time.Location {
	return t.TLoc
}

func (t *Trigger) NextTriggerTime() time.Time {
	return t.TNextTime
}

func ModifyTrigger(t triggers.ImmutableTrigger, f func(tr *Trigger)) triggers.ImmutableTrigger {
	if trigger, ok := t.(*Trigger); ok {
		f(trigger)
	}
	return t
}

func NewTrigger() *Trigger {
	return &Trigger{
		TRepeats:  triggers.RepeatInfinity,
		TLocation: "Local",
	}
}
