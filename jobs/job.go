package jobs

import "errors"

var (
	ErrEmptyJobKey  = errors.New("empty job key")
	ErrEmptyJobType = errors.New("empty job type")
)

type MutableJob interface {
	WithData(data interface{}) MutableJob
	WithKey(jKey string) MutableJob
	WithType(jType string) MutableJob
	ToImmutable() (ImmutableJob, error)
}

type ImmutableJob interface {
	Key() string
	Type() string
	Data() []byte
}
