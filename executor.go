package scheduler

import (
	"github.com/d1slike/go-sched/log"
	"github.com/d1slike/go-sched/stores"
	"sync"
	"time"
)

type executor interface {
	Start()
	Shutdown() error
	CancelTrigger(tKey string) bool
}

type defaultRuntimeExecutor struct {
	sName    string
	store    stores.Store
	registry executorRegistry
	timers   Timers

	tLock sync.Mutex
	tMap  map[string]*time.Timer

	closeChan chan struct{}
}

func (e *defaultRuntimeExecutor) CancelTrigger(tKey string) bool {
	e.tLock.Lock()
	defer e.tLock.Unlock()

	t, ok := e.tMap[tKey]
	if ok {
		t.Stop()
		delete(e.tMap, tKey)
	}

	return ok
}

func (e *defaultRuntimeExecutor) Start() {
	e.startTriggerStealing()
}

func (e *defaultRuntimeExecutor) Shutdown() error {
	close(e.closeChan)
	return nil
}

func (e *defaultRuntimeExecutor) startTriggerStealing() {
	go func() {
		for {
			select {
			case <-e.closeChan:
				return
			case <-time.After(e.timers.TriggerStealTimeout):
				log.Debug("defaultRuntimeExecutor: stealing free triggers")
			}
		}
	}()
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
