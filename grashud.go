package grashud

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

type shutdownFunc = func() error

type grashud struct {
	funcs   []shutdownFunc
	once    sync.Once
	errCh   chan error
	sigCh   chan os.Signal
	isPanic chan struct{}
}

func New() *grashud {
	g := &grashud{
		funcs:   make([]shutdownFunc, 0),
		once:    sync.Once{},
		isPanic: make(chan struct{}),
		errCh:   make(chan error),
		sigCh:   make(chan os.Signal),
	}

	return g
}

func (g *grashud) Add(funcs ...shutdownFunc) {
	g.funcs = append(g.funcs, funcs...)
}

func (g *grashud) HandleSignals() {
	go g.once.Do(func() {
		logger, _ := zap.NewProduction()
		signal.Notify(g.sigCh, syscall.SIGINT, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGSEGV)
		select {
		case <-g.sigCh:
			logger.Error("Signal recived")
			g.callFuncs()
		case <-g.isPanic:
			logger.Error("Panic arised")
			g.callFuncs()
		}
		for e := range g.errCh {
			logger.Error(e.Error())
			logger.Sync()
		}
	})
}

func (g *grashud) HandlePanic() {
	if r := recover(); r != nil {
		g.isPanic <- struct{}{}
		panic(r)
	}
}

func (g *grashud) callFuncs() {
	for _, f := range g.funcs {
		go func(f shutdownFunc) {
			err := f()
			if err != nil {
				g.errCh <- err
			}
		}(f)
	}
}
