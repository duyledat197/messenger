package processor

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewLifecycle returns an instance of the Lifecycle.
func NewLifecycle() *Lifecycle {
	return &Lifecycle{}
}

// WithProcessors appends the given processors to the Lifecycle's processors slice.
func (l *Lifecycle) WithProcessors(processors ...Processor) {
	l.processors = append(l.processors, processors...)
}

// WithFactories appends the given factories to the Lifecycle's factories slice.
func (l *Lifecycle) WithFactories(factories ...Factory) {
	l.factories = append(l.factories, factories...)
}

// Start starts the lifecycle with the given context.
func (l *Lifecycle) Start(ctx context.Context) {
	errChan := make(chan error)
	signChan := make(chan os.Signal, 1)

	for _, f := range l.factories {
		if err := f.Connect(ctx); err != nil {
			slog.Error("unable to connect factory:", slog.String("err", err.Error())) //nolint:lll //err)
		}
	}

	for _, p := range l.processors {
		go func(p Processor) {
			if err := p.Start(ctx); err != nil {
				errChan <- err
			}
		}(p)
	}

	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-errChan:
		l.stop(ctx)
	case <-signChan:
		l.stop(ctx, true)
	}

}

// stop stops the server gracefully by closing all factories and starting all processors.
func (l *Lifecycle) stop(ctx context.Context, graceful ...bool) {
	for _, p := range l.processors {
		if err := p.Stop(ctx); err != nil {
			slog.Error("unable to close processor:", slog.String("err", err.Error()))
		}
	}
	slog.Info("graceful shutdown...")
	if len(graceful) > 0 {
		time.Sleep(5 * time.Second)
	}

	for _, f := range l.factories {
		if err := f.Close(ctx); err != nil {
			slog.Error("unable to close factory:", slog.String("err", err.Error()))
		}
	}
}
