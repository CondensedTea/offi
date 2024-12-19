package closer

import (
	"context"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type funcMock struct {
	completedCallbacks atomic.Int64
	sleep              time.Duration
}

func (m *funcMock) callback() error {
	defer m.completedCallbacks.Add(1)

	time.Sleep(m.sleep)

	return nil
}

func Test_CloseAll(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		// setup
		const (
			callbackRunningDuration = 15 * time.Millisecond
			wantCallbacks           = 2
		)

		var m = funcMock{
			sleep: callbackRunningDuration,
		}

		var c = Closer{
			logger: slog.Default(),
		}

		for i := 0; i < wantCallbacks; i++ {
			c.Add(m.callback)
		}

		now := time.Now()

		// act
		c.CloseAll(context.Background())

		// assert
		assert.Equal(t, int64(wantCallbacks), m.completedCallbacks.Load())
		assert.InEpsilon(t, callbackRunningDuration, time.Since(now), 0.08)
	})

	t.Run("context canceled", func(t *testing.T) {
		t.Parallel()

		// setup
		var (
			fastMock = funcMock{sleep: time.Millisecond}
			slowMock = funcMock{sleep: 10 * time.Second}
		)

		var c = Closer{
			logger: slog.Default(),
		}

		c.Add(fastMock.callback)
		c.Add(slowMock.callback)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		defer cancel()

		// act
		c.CloseAll(ctx)

		// assert
		assert.Equal(t, int64(1), fastMock.completedCallbacks.Load())
		assert.Equal(t, int64(0), slowMock.completedCallbacks.Load())
	})
}
