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
	lock           sync.Mutex
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

	if cce.errorCollector == nil {
		cce.errorCollector = err
	} else {
		cce.errorCollector = fmt.Errorf("%s; %s", cce.errorCollector, err.Error())
	}
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
	return cce.runningNum
}

func (cce *ConControllerWithError) Wait(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("concurrent controller cancel wait for context canceled")
			return

		case <-ticker.C:
			if cce.runningNum == 0 {
				return
			}
		}
	}
}
