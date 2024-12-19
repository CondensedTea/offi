package closer

import (
	"context"
	"log/slog"
	"sync"
)

type closeFunc func() error

type closeContextFunc func(ctx context.Context) error

var global *Closer

// Add adds callback to global Closer.
func Add(f closeFunc) {
	global.Add(f)
}

// AddContext adds callback that uses context to global Closer.
func AddContext(f closeContextFunc) {
	global.AddContext(f)
}

// CloseAll runs callbacks stored in global Closer.
func CloseAll(ctx context.Context) {
	global.CloseAll(ctx)
}

func init() {
	global = &Closer{
		logger: slog.Default().With("component", "closer"),
	}
}

// Closer allows to collect shutodown callbacks and execute them on request.
type Closer struct {
	logger *slog.Logger
	funcs  []closeContextFunc
}

// Add stores callback in Closer
func (c *Closer) Add(f closeFunc) {
	c.funcs = append(c.funcs, func(context.Context) error {
		return f()
	})
}

// AddContext adds callback that uses context to global Closer.
func (c *Closer) AddContext(f closeContextFunc) {
	c.funcs = append(c.funcs, f)
}

// CloseAll runs all stored callbacks in parallel and waits for them or until provided context is canceled.
func (c *Closer) CloseAll(ctx context.Context) {
	wg := sync.WaitGroup{}

	for _, f := range c.funcs {
		callback := f

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := callback(ctx); err != nil {
				c.logger.Error("failed to call closer func", "error", err)
			}
		}()
	}

	var completeCh = make(chan struct{})

	go func() {
		wg.Wait()
		completeCh <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			c.logger.Warn("some closers took too long to finish", "component", "closer")
			return
		case <-completeCh:
			return
		}
	}
}
