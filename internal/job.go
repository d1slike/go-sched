package internal

import (
	"github.com/d1slike/go-sched/jobs"
	"github.com/d1slike/go-sched/log"
)

type Job struct {
	Jkey   string
	JjType string
	Jdata  []byte
}

func (j *Job) ToImmutable() (jobs.ImmutableJob, error) {
	if j.Jkey == "" {
		return nil, jobs.ErrEmptyJobKey
	}
	if j.JjType == "" {
		return nil, jobs.ErrEmptyJobType
	}
	return j, nil
}

func (j *Job) Key() string {
	return j.Jkey
}

func (j *Job) Type() string {
	return j.JjType
}

func (j *Job) Data() []byte {
	return j.Jdata
}

func (j *Job) WithData(data interface{}) jobs.MutableJob {
	if b, err := castData(data); err != nil {
		log.Errorf("job key: %s, jon type: %s : %v", j.Jkey, j.JjType, err)
	} else {
		j.Jdata = b
	}
	return j
}

func (j *Job) WithKey(jKey string) jobs.MutableJob {
	j.Jkey = jKey
	return j
}

func (j *Job) WithType(jType string) jobs.MutableJob {
	j.JjType = jType
	return j
}

func NewJob() *Job {
	return &Job{}
}
