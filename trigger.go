package scheduler

import (
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"time"
)

const (
	Infinity = Repeats(-1)
)

var (
	ErrEmptyTriggerKey = errors.New("empty trigger key")
	ErrEmptyCronSpec   = errors.New("empty cron specification")
	ErrInvalidLocation = "invalid location: %v"
	ErrInvalidCronSpec = "invalid cron spec: %v"
)

type Repeats int

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
	NextTriggerTime() time.Time
}

type trigger struct {
	key      string
	jobKey   string
	fromTime *time.Time
	toTime   *time.Time
	repeats  Repeats
	cronSpec string
	location string

	loc      *time.Location
	sched    cron.Schedule
	nextTime time.Time
}

func (t *trigger) JobKey() string {
	return t.key
}

func (t *trigger) WithKey(tKey string) MutableTrigger {
	t.key = tKey
	return t
}

func (t *trigger) WithFromTime(from time.Time) MutableTrigger {
	t.fromTime = &from
	return t
}

func (t *trigger) WithToTime(to time.Time) MutableTrigger {
	t.toTime = &to
	return t
}

func (t *trigger) WithRepeats(repeats Repeats) MutableTrigger {
	t.repeats = repeats
	return t
}

func (t *trigger) WithCron(spec string) MutableTrigger {
	t.cronSpec = spec
	return t
}

func (t *trigger) InLocation(loc string) MutableTrigger {
	t.location = loc
	return t
}

func (t *trigger) ToImmutable() (ImmutableTrigger, error) {
	if t.key == "" {
		return nil, ErrEmptyTriggerKey
	}
	if t.cronSpec == "" {
		return nil, ErrEmptyCronSpec
	}

	loc, err := time.LoadLocation(t.location)
	if err != nil {
		return nil, fmt.Errorf(ErrInvalidLocation, err)
	}
	t.loc = loc

	sched, err := cron.Parse(t.cronSpec)
	if err != nil {
		return nil, fmt.Errorf(ErrInvalidCronSpec, err)
	}
	t.sched = sched

	t.nextTime = t.sched.Next(time.Now().In(t.loc))

	return t, nil
}

func (t *trigger) Key() string {
	return t.key
}

func (t *trigger) FromTime() *time.Time {
	return t.fromTime
}

func (t *trigger) ToTime() *time.Time {
	return t.toTime
}

func (t *trigger) Repeats() Repeats {
	return t.repeats
}

func (t *trigger) CronSpec() string {
	return t.cronSpec
}

func (t *trigger) Location() *time.Location {
	return t.loc
}

func (t *trigger) NextTriggerTime() time.Time {
	return t.nextTime
}

func NewTrigger() MutableTrigger {
	return &trigger{
		repeats:  Infinity,
		location: "Local",
	}
}

func Repeat(count int) Repeats {
	return Repeats(count)
}
