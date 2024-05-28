package concurrent

import (
	"context"
	"sync"
	"time"
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
	if cc.allowSize == 0 {
		return
	}

	ticker := time.NewTicker(time.Millisecond * 10)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		if cc.RunningNum() == 0 {
			return
		}
	}
}

func (cc *ConController) Close() {
	if cc.workerChan != nil {
		close(cc.workerChan)
	}
}
