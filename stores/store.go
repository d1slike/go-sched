package stores

import (
	"errors"
	"github.com/d1slike/go-sched/jobs"
	"github.com/d1slike/go-sched/triggers"
)

var (
	ErrJobAlreadyExists     = errors.New("job with same key already exists")
	ErrTriggerAlreadyExists = errors.New("trigger with same key already exists")
	ErrJobNotFound          = errors.New("job not found")
	ErrTriggerNotFound      = errors.New("trigger not found")
)

type Store interface {
	InsertJob(sName string, job jobs.ImmutableJob) error
	InsertTrigger(sName string, trigger triggers.ImmutableTrigger) error
	GetJob(sName string, jKey string) (jobs.ImmutableJob, error)
	GetTrigger(sName string, tKey string) (triggers.ImmutableTrigger, error)
	DeleteJob(sName string, jKey string) (bool, error)
	DeleteTrigger(sName string, tKey string) (bool, error)
	DeleteTriggersByJobKey(sName string, jKey string) ([]string, error)
	GetJobs(sName string) ([]jobs.ImmutableJob, error)
	GetTriggers(sName string) ([]triggers.ImmutableTrigger, error)
	AcquireTriggers(sName string) ([]triggers.ImmutableTrigger, error)
	UpdateTrigger(sName string, trigger triggers.ImmutableTrigger) error
	UpdateJob(sName string, job jobs.ImmutableJob) error
	DeleteExhaustedTriggers(sName string) (int, error)
}
