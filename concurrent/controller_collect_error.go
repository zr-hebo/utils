package concurrent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ConControllerWithError struct {
	allowSize      int
	runningNum     int
	lock           sync.RWMutex
	workerChan     chan struct{}
	errorCollector error
}

func NewConControllerWithError(size int) (cc *ConControllerWithError) {
	return &ConControllerWithError{
		allowSize:  size,
		workerChan: make(chan struct{}, size),
	}
}

func (cce *ConControllerWithError) Acquire() {
	cce.workerChan <- struct{}{}
	cce.lock.Lock()
	cce.runningNum++
	cce.lock.Unlock()
}

func (cce *ConControllerWithError) CollectError(err error) {
	if err == nil {
		return
	}

	cce.lock.Lock()
	if cce.errorCollector == nil {
		cce.errorCollector = err
	} else {
		cce.errorCollector = fmt.Errorf("%s; %s", cce.errorCollector, err.Error())
	}
	cce.lock.Unlock()
}

func (cce *ConControllerWithError) Error() error {
	return cce.errorCollector
}

func (cce *ConControllerWithError) Release() {
	<-cce.workerChan
	cce.lock.Lock()
	cce.runningNum--
	cce.lock.Unlock()
}

func (cce *ConControllerWithError) RunningNum() int {
	cce.lock.RLock()
	defer cce.lock.RUnlock()
	return cce.runningNum
}

func (cce *ConControllerWithError) Wait(ctx context.Context) {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("concurrent controller cancel wait for context canceled")
			return

		case <-ticker.C:
			if cce.RunningNum() == 0 {
				return
			}
		}
	}
}
