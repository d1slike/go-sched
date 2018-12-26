package internal

import (
	"github.com/d1slike/go-sched/jobs"
	"github.com/d1slike/go-sched/json"
	"github.com/d1slike/go-sched/log"
)

type Job struct {
	JKey   string
	JJType string
	JData  []byte
}

func (j *Job) ToImmutable() (jobs.ImmutableJob, error) {
	if j.JKey == "" {
		return nil, jobs.ErrEmptyJobKey
	}
	if j.JJType == "" {
		return nil, jobs.ErrEmptyJobType
	}
	return j, nil
}

func (j *Job) Key() string {
	return j.JKey
}

func (j *Job) Type() string {
	return j.JJType
}

func (j *Job) Data() []byte {
	return j.JData
}

func (j *Job) WithData(data interface{}) jobs.MutableJob {
	switch d := data.(type) {
	case []byte:
		j.JData = d
	case string:
		j.JData = []byte(d)
	default:
		if bytes, err := json.Provider.Marshal(data); err != nil {
			log.Errorf("job key: %s, jon type: %s : %v", j.JKey, j.JJType, err)
		} else {
			j.JData = bytes
		}
	}

	return j
}

func (j *Job) WithKey(jKey string) jobs.MutableJob {
	j.JKey = jKey
	return j
}

func (j *Job) WithType(jType string) jobs.MutableJob {
	j.JJType = jType
	return j
}

func NewJob() *Job {
	return &Job{}
}
