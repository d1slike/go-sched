package scheduler

import "time"

type Repeats int

const (
	Infinity = Repeats(-1)
)

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
	FromTime() *time.Time
	ToTime() *time.Time
	Repeats() Repeats
	CronSpec() string
	Location() *time.Location
	NextTriggerTime() time.Time
}

type trigger struct {
	key      string
	fromTime *time.Time
	toTime   *time.Time
	repeats  Repeats
	cronSpec string
	location string

	loc      *time.Location
	nextTime time.Time
}

func (t *trigger) WithKey(tKey string) MutableTrigger {
	return t
}

func (t *trigger) WithFromTime(from time.Time) MutableTrigger {
	return t
}

func (t *trigger) WithToTime(to time.Time) MutableTrigger {
	return t
}

func (t *trigger) WithRepeats(repeats Repeats) MutableTrigger {
	return t
}

func (t *trigger) WithCron(spec string) MutableTrigger {
	return t
}

func (t *trigger) InLocation(loc string) MutableTrigger {
	return t
}

func (t *trigger) ToImmutable() (ImmutableTrigger, error) {
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
	return &trigger{}
}

func Repeat(count int) Repeats {
	return Repeats(count)
}
