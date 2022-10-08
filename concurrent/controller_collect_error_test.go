package concurrent

import (
	"context"
	"testing"
)

func TestRace(t *testing.T) {
	n := 10000
	cc := NewConControllerWithError(n)
	for i := 0; i < n; i++ {
		cc.Acquire()

		go func() {
			cc.Release()
		}()
	}

	cc.Wait(context.Background())
}
