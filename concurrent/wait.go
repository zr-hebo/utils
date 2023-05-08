package concurrent

import (
	"context"
	"time"
)

func WaitWithDuration(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		return
	}
}
