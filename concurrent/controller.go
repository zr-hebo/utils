package concurrent

import (
	"context"
	"fmt"
)

type ConController struct {
	allowSize  int
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
}

func (cc *ConController) Release() {
	<-cc.workerChan
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
