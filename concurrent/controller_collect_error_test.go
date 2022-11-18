package concurrent

import (
	"context"
	"testing"
)

func TestRace(t *testing.T) {
	execOnceCheck()
}

func execOnceCheck() {
	n := 10000
	cc := NewConControllerWithError(100)
	for i := 0; i < n; i++ {
		cc.Acquire()

		go func() {
			cc.Release()
		}()
	}

	cc.Wait(context.Background())
}

func BenchmarkConController(b *testing.B) {
	for i := 0; i < b.N; i++ {
		execOnceCheck()
	}
}
