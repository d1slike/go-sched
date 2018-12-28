package scheduler

import (
	"context"
	"errors"
	"fmt"
	"github.com/d1slike/go-sched/internal"
	"github.com/d1slike/go-sched/log"
	"github.com/d1slike/go-sched/stores"
	"github.com/d1slike/go-sched/triggers"
	"github.com/d1slike/go-sched/utils"
	"sync"
	"time"
)

const (
	defaultDelta = 10 * time.Second
)

var (
	ErrJobDeadlineExceeded = errors.New("job execution deadline exceeded")
)

type executor interface {
	Start()
	Shutdown(ctx context.Context) error
	CancelTriggers(tKey ...string) int
}

type defaultRuntimeExecutor struct {
	sName    string
	store    stores.Store
	registry executorRegistry
	timers   Timers

	runningFutures sync.WaitGroup
	lock           sync.Mutex
	fMap           map[string]*future

	closeChan chan struct{}
}

func (e *defaultRuntimeExecutor) CancelTriggers(keys ...string) int {
	e.lock.Lock()
	defer e.lock.Unlock()

	canceled := 0
	for _, tKey := range keys {
		f, ok := e.fMap[tKey]
		if ok {
			f.Cancel()
			delete(e.fMap, tKey)
			canceled++
		}
	}

	return canceled
}

func (e *defaultRuntimeExecutor) Start() {
	e.startTriggerStealing()
}

func (e *defaultRuntimeExecutor) Shutdown(ctx context.Context) error {
	//stop all internal background tasks
	close(e.closeChan)

	//try release not running triggers
	e.lock.Lock()
	for key, f := range e.fMap {
		if !f.IsRunning() {
			f.Cancel()
			f.t = internal.ModifyTrigger(f.t, func(tr *internal.Trigger) {
				tr.Tstate = triggers.StateScheduled //just release scheduled triggers
			})
			if err := e.store.UpdateTrigger(e.sName, f.t); err != nil {
				log.Errorf("defaultRuntimeExecutor: could not update trigger %v: %v", f.t.Key(), err)
			}
			delete(e.fMap, key)
		}
	}
	e.lock.Unlock()

	//await all running triggers
	awaitRunning := make(chan struct{}, 1)
	go func() {
		e.runningFutures.Wait()
		awaitRunning <- struct{}{}
	}()
	select {
	case <-ctx.Done():
	case <-awaitRunning:

	}

	//release remaining triggers
	e.lock.Lock()
	/*for key, f := range e.fMap {

	}*/
	e.lock.Unlock()

	return nil
}

func (e *defaultRuntimeExecutor) startTriggerStealing() {
	go func() {
		for {
			select {
			case <-e.closeChan:
				return
			case <-time.After(e.timers.TriggerStealTimeout):
				triggers, err := e.store.AcquireTriggers(e.sName)

				if err != nil {
					log.Errorf("defaultRuntimeExecutor: could not acquire free triggers: %v", err)
				} else {
					log.Debugf("defaultRuntimeExecutor: acquired %d free triggers", len(triggers))

					if len(triggers) > 0 {
						e.lock.Lock()
						for _, t := range triggers {
							e.fMap[t.Key()] = e.makeFuture(t)
						}
						e.lock.Unlock()
					}
				}

			}
		}
	}()
}

func (e *defaultRuntimeExecutor) makeFuture(t triggers.ImmutableTrigger) *future {
	future := &future{
		t:        t,
		running:  utils.NewAtomicBool(false),
		canceled: utils.NewAtomicBool(false),
	}
	f := e.makeF(future)

	var dur time.Duration
	now := time.Now().In(t.Location())
	if now.After(t.NextTriggerTime()) {
		dur = 0
	} else {
		dur = t.NextTriggerTime().Sub(now)
	}

	future.timer = time.AfterFunc(dur, f)

	return future
}

func (e *defaultRuntimeExecutor) makeF(f *future) func() {
	return func() {
		if f.IsCanceled() {
			return
		}

		f.Run()
		e.runningFutures.Add(1)
		defer func() {
			e.lock.Lock()
			delete(e.fMap, f.t.Key())
			e.lock.Unlock()

			f.running.Set(false)

			e.runningFutures.Done()
		}()

		trigger, err := e.store.GetTrigger(e.sName, f.t.Key())
		if err != nil {
			log.Errorf("defaultRuntimeExecutor: could not get trigger %v: %v", f.t.Key(), err)
			return
		}
		if trigger == nil {
			log.Warnf("defaultRuntimeExecutor: trigger %v was deleted", f.t.Key())
			return
		}
		now := time.Now().In(trigger.Location())
		if !internal.IsNear(now, trigger.NextTriggerTime(), defaultDelta) {
			log.Warnf("defaultRuntimeExecutor: trigger %v was updated. now: %s, next trigger time: %s", f.t.Key(), now, trigger.NextTriggerTime())
			return
		}

		job, err := e.store.GetJob(e.sName, trigger.JobKey())
		if err != nil {
			log.Errorf("defaultRuntimeExecutor: could not get job %v: %v", trigger.JobKey(), err)
			return
		}
		if job == nil {
			log.Warnf("defaultRuntimeExecutor: job %v was deleted", trigger.JobKey())
			return
		}

		exec, ok := e.registry.GetExecutor(job.Type())
		if !ok {
			log.Errorf("defaultRuntimeExecutor: not found executor for job type: %v", job.Type())
			return
		}

		ctx := &jobCtx{
			job:     job,
			trigger: trigger,
		}
		doneChan := make(chan error, 1)
		timeout := make(chan time.Time) //todo add timeout
		go func() {
			defer func() {
				if err := recover(); err != nil {
					doneChan <- fmt.Errorf("%v", err)
				}
			}()
			doneChan <- exec(ctx)
		}()

		select {
		case <-timeout:
			err = ErrJobDeadlineExceeded
		case e := <-doneChan:
			err = e
		}

		if err != nil {
			log.Warnf("defaultRuntimeExecutor: job %v has finished with err '%v' by trigger %v", job.Key(), err, trigger.Key())
		}

		trigger = internal.ModifyTrigger(trigger, func(tr *internal.Trigger) {
			tr.TtriggeredTime++
			if tr.Trepeats != triggers.RepeatInfinity && tr.TtriggeredTime >= tr.Trepeats {
				tr.Tstate = triggers.StateExhausted
				return
			}

			nextTime := internal.CalcNextTriggerTime(tr)
			if nextTime.IsZero() {
				tr.Tstate = triggers.StateExhausted
			} else {
				tr.Tstate = triggers.StateScheduled
				tr.TnextTime = nextTime
			}
		})

		if err := e.store.UpdateTrigger(e.sName, trigger); err != nil {
			log.Errorf("defaultRuntimeExecutor: could not update trigger %v: %v", trigger.Key(), err)
		}
	}
}

func newDefaultRuntimeExecutor(
	sName string,
	store stores.Store,
	registry executorRegistry,
	timers Timers,
) executor {
	return &defaultRuntimeExecutor{
		sName:     sName,
		store:     store,
		registry:  registry,
		timers:    timers,
		closeChan: make(chan struct{}),
		fMap:      make(map[string]*future),
	}
}
