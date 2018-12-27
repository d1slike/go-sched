package internal

import (
	"fmt"
	"github.com/d1slike/go-sched/log"
	"github.com/d1slike/go-sched/triggers"
	"github.com/robfig/cron"
	"time"
)

type Trigger struct {
	Tkey      string
	TjobKey   string
	TfromTime *time.Time
	TtoTime   *time.Time
	Trepeats  triggers.Repeats
	TcronSpec string
	Tlocation string
	Tdata     []byte

	Tstate         triggers.TriggerState
	Tloc           *time.Location
	TtriggeredTime triggers.Repeats
	Tsched         cron.Schedule
	TnextTime      time.Time
}

func (t *Trigger) Data() []byte {
	return t.Tdata
}

func (t *Trigger) TriggeredTimes() triggers.Repeats {
	return t.TtriggeredTime
}

func (t *Trigger) State() triggers.TriggerState {
	return t.Tstate
}

func (t *Trigger) JobKey() string {
	return t.Tkey
}

func (t *Trigger) WithKey(tKey string) triggers.MutableTrigger {
	t.Tkey = tKey
	return t
}

func (t *Trigger) WithFromTime(from time.Time) triggers.MutableTrigger {
	t.TfromTime = &from
	return t
}

func (t *Trigger) WithToTime(to time.Time) triggers.MutableTrigger {
	t.TtoTime = &to
	return t
}

func (t *Trigger) WithRepeats(repeats triggers.Repeats) triggers.MutableTrigger {
	t.Trepeats = repeats
	return t
}

func (t *Trigger) WithCron(spec string) triggers.MutableTrigger {
	t.TcronSpec = spec
	return t
}

func (t *Trigger) WithData(data interface{}) triggers.MutableTrigger {
	if b, err := castData(data); err != nil {
		log.Errorf("trigger key: %v", t.Tkey, err)
	} else {
		t.Tdata = b
	}
	return t
}

func (t *Trigger) InLocation(loc string) triggers.MutableTrigger {
	t.Tlocation = loc
	return t
}

func (t *Trigger) ToImmutable() (triggers.ImmutableTrigger, error) {
	if t.Tkey == "" {
		return nil, triggers.ErrEmptyTriggerKey
	}
	if t.TcronSpec == "" {
		return nil, triggers.ErrEmptyCronSpec
	}

	loc, err := time.LoadLocation(t.Tlocation)
	if err != nil {
		return nil, fmt.Errorf(triggers.ErrInvalidLocation, err)
	}
	t.Tloc = loc

	sched, err := cron.Parse(t.TcronSpec)
	if err != nil {
		return nil, fmt.Errorf(triggers.ErrInvalidCronSpec, err)
	}
	t.Tsched = sched

	nextTime := CalcNextTriggerTime(t)
	if nextTime.IsZero() {
		return nil, triggers.ErrAlreadyExhausted
	} else {
		t.TnextTime = nextTime
	}

	return t, nil
}

func (t *Trigger) Key() string {
	return t.Tkey
}

func (t *Trigger) FromTime() *time.Time {
	return t.TfromTime
}

func (t *Trigger) ToTime() *time.Time {
	return t.TtoTime
}

func (t *Trigger) Repeats() triggers.Repeats {
	return t.Trepeats
}

func (t *Trigger) CronSpec() string {
	return t.TcronSpec
}

func (t *Trigger) Location() *time.Location {
	return t.Tloc
}

func (t *Trigger) NextTriggerTime() time.Time {
	return t.TnextTime
}

func ModifyTrigger(t triggers.ImmutableTrigger, f func(tr *Trigger)) triggers.ImmutableTrigger {
	if trigger, ok := t.(*Trigger); ok {
		cpy := *trigger
		f(&cpy)
		return &cpy
	}
	return t
}

func NewTrigger() *Trigger {
	return &Trigger{
		Trepeats:  triggers.RepeatInfinity,
		Tlocation: "Local",
	}
}

// calc next trigger time considering fromTime, toTime boundary
// return zero time if never fire
func CalcNextTriggerTime(t *Trigger) time.Time {
	var nextTime time.Time
	if t.TfromTime != nil {
		nextTime = t.Tsched.Next(t.TfromTime.In(t.Tloc))
	} else {
		nextTime = t.Tsched.Next(time.Now().In(t.Tloc))
	}

	if t.TnextTime.IsZero() || (t.TtoTime != nil && t.TtoTime.Before(t.TnextTime)) {
		return time.Time{}
	}

	return nextTime
}

func IsNear(a, b time.Time, delta time.Duration) bool {
	if a.Equal(b) {
		return true
	}

	if a.After(b) {
		return a.Sub(b) <= delta
	}

	return b.Sub(a) <= delta
}
