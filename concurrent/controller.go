package concurrent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ConController struct {
	allowSize  int
	runningNum int
	lock       sync.Mutex
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
	<-cc.workerChan
	cc.lock.Lock()
	cc.runningNum--
	cc.lock.Unlock()
}

func (cc *ConController) RunningNum() int {
	return cc.runningNum
}

func (cc *ConController) Wait(ctx context.Context) {
	ticker := time.NewTicker(time.Millisecond * 10)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("concurrent controller cancel wait for context canceled")
			return

		case <-ticker.C:
			if cc.runningNum == 0 {
				return
			}
		}
	}
	return
}
