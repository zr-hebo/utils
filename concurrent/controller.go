package concurrent

import (
	"context"
	"sync"
)

type ConController struct {
	allowSize  int
	runningNum int
	lock       sync.RWMutex
	workerChan chan struct{}
}

func NewConController(size int) (cc *ConController) {
	return &ConController{
		allowSize:  size,
		workerChan: make(chan struct{}, size),
	}
}

func (cc *ConController) Acquire() {
	cc.workerChan <- struct{}{}
	cc.lock.Lock()
	cc.runningNum++
	cc.lock.Unlock()
}

func (cc *ConController) Release() {
	cc.lock.Lock()
	cc.runningNum--
	cc.lock.Unlock()
	<-cc.workerChan
}

func (cc *ConController) RunningNum() int {
	cc.lock.RLock()
	defer cc.lock.RUnlock()
	return cc.runningNum
}

func (cc *ConController) Wait(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(cc.workerChan)
			return
		case cc.workerChan <- struct{}{}:
		}

		if cc.RunningNum() == 0 {
			close(cc.workerChan)
			return
		}
	}
}
