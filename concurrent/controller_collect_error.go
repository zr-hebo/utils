package concurrent

import (
	"context"
	"fmt"
	"sync"
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
	cce.lock.RLock()
	defer cce.lock.RUnlock()
	return cce.errorCollector
}

func (cce *ConControllerWithError) Release() {
	cce.lock.Lock()
	cce.runningNum--
	cce.lock.Unlock()
	<-cce.workerChan
}

func (cce *ConControllerWithError) RunningNum() int {
	cce.lock.RLock()
	defer cce.lock.RUnlock()
	return cce.runningNum
}

func (cce *ConControllerWithError) Size() int {
	return cce.allowSize
}

func (cce *ConControllerWithError) Wait(ctx context.Context) {
	defer func() {
		close(cce.workerChan)
	}()
	if cce.allowSize == 0 {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case cce.workerChan <- struct{}{}:
		}

		if cce.RunningNum() == 0 {
			return
		}
	}
}
