package scheduler

import (
	"context"
	"github.com/d1slike/go-sched/internal"
	"github.com/d1slike/go-sched/log"
	"github.com/d1slike/go-sched/stores"
	"github.com/d1slike/go-sched/triggers"
	"github.com/d1slike/go-sched/utils"
	"sync"
	"time"
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
	awaitRun := make(chan struct{}, 1)
	go func() {
		e.runningFutures.Wait()
		awaitRun <- struct{}{}
	}()
	select {
	case <-ctx.Done():
	case <-awaitRun:

	}

	//release remaining triggers
	e.lock.Lock()
	for key, f := range e.fMap {

	}
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
	f := func() {

	}

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
	}
}
