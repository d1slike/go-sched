package scheduler

import (
	"errors"
	"github.com/d1slike/go-sched/json"
	"github.com/d1slike/go-sched/log"
)

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

type job struct {
	key   string
	jType string
	data  []byte
}

func (j *job) ToImmutable() (ImmutableJob, error) {
	if j.key == "" {
		return nil, ErrEmptyJobKey
	}
	if j.jType == "" {
		return nil, ErrEmptyJobType
	}
	return j, nil
}

func (j *job) Key() string {
	return j.key
}

func (j *job) Type() string {
	return j.jType
}

func (j *job) Data() []byte {
	return j.data
}

func (j *job) WithData(data interface{}) MutableJob {
	switch d := data.(type) {
	case []byte:
		j.data = d
	case string:
		j.data = []byte(d)
	default:
		if bytes, err := json.Provider.Marshal(data); err != nil {
			log.Errorf("job key: %s, jon type: %s : %v", j.key, j.jType, err)
		} else {
			j.data = bytes
		}
	}

	return j
}

func (j *job) WithKey(jKey string) MutableJob {
	j.key = jKey
	return j
}

func (j *job) WithType(jType string) MutableJob {
	j.jType = jType
	return j
}

func NewJob() MutableJob {
	return &job{}
}
