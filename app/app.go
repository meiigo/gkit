package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/meiigo/gkit/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// App is an application components lifecycle manager
type App struct {
	opts   options
	ctx    context.Context
	cancel func()
}

// New create an application lifecycle manager.
func New(opts ...Option) *App {
	o := options{
		ctx:         context.Background(),
		logger:      log.NewHelper(log.DefaultLogger),
		sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		stopTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.opts.id }

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Metadata returns service metadata.
func (a *App) Metadata() map[string]string { return a.opts.metadata }

type appKey struct{}

func (a *App) Run() error {
	eg, ctx := errgroup.WithContext(a.ctx)
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			stopCtx, cancel := context.WithTimeout(context.WithValue(a.opts.ctx, appKey{}, a), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		eg.Go(func() error {
			return srv.Start(context.WithValue(a.opts.ctx, appKey{}, a))
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				if err := a.Stop(); err != nil {
					a.opts.logger.Errorf("failed to stop app: %v", err)
					return err
				}
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
