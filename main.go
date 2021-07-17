package grashud

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type shutdownFunc = func() error
type edge int

const (
	WithCancel edge = iota
	WithDeadline
	WithTimeout
)

type grashud struct {
	funcs []shutdownFunc
	once  sync.Once
	ctx   context.Context
}

func New(e edge, v interface{}) (*grashud, context.Context, context.CancelFunc) {
	base := context.Background()
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	switch e {
	case WithCancel:
		ctx, cancel = context.WithCancel(base)
	case WithDeadline:
		ctx, cancel = context.WithDeadline(base, v.(time.Time))
	case WithTimeout:
		ctx, cancel = context.WithTimeout(base, v.(time.Duration))
	}

	g := &grashud{
		funcs: make([]shutdownFunc, 0),
		once:  sync.Once{},
		ctx:   ctx,
	}

	go g.start()

	return g, ctx, cancel
}

func (g *grashud) AddFunc(funcs ...shutdownFunc) {
	g.funcs = append(g.funcs, funcs...)
}

func (g *grashud) start() {
	g.once.Do(func() {
		logger, _ := zap.NewProduction()
		sigCh, errCh := make(chan os.Signal, 1), make(chan error)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP)
		select {
		case <-sigCh:
			logger.Info("Signal recived, trying to graceful shutdown")
			for _, f := range g.funcs {
				go func(f shutdownFunc) {
					err := f()
					if err != nil {
						errCh <- err
					}
				}(f)
			}
			close(sigCh)
		case <-g.ctx.Done():
			logger.Info("Graceful shutdown")
			for _, f := range g.funcs {
				go func(f shutdownFunc) {
					err := f()
					if err != nil {
						errCh <- err
					}
				}(f)
			}
			close(sigCh)
		}
		for e := range errCh {
			logger.Error(e.Error())
			logger.Sync()
		}
	})
}
