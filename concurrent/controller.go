package concurrent

import (
    "context"
    "fmt"
    "sync"
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

func (cc *ConController) Wait(ctx context.Context) (err error) {
    for idx := 0; idx < cc.allowSize; idx++ {
        select {
        case <-ctx.Done():
            err = fmt.Errorf("concurrent controller cancel wait for context calceled")
            return
        case cc.workerChan <- struct{}{}:
        }
    }
    return
}
