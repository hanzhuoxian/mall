package shutdown

import (
	"sync"
)

type ShutdownCallback interface {
	OnShutdown(string) error
}

type ShutdownFunc func(string) error

func (f ShutdownFunc) OnShutdown(s string) error {
	return f(s)
}

type ErrorHandler interface {
	OnError(err error)
}

type ErrorFunc func(err error)

func (f ErrorFunc) OnError(err error) {
	f(err)
}

type ShutdownManager interface {
	GetName() string
	Start(gs GracefulShutter) error
	ShutdownStart() error
	ShutdownFinish() error
}

type GracefulShutter interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(sc ShutdownCallback)
}

type GracefulShutdown struct {
	callbacks        []ShutdownCallback
	managers         []ShutdownManager
	errorHandler     ErrorHandler
	shutdownFinished chan struct{}
}

func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks:        make([]ShutdownCallback, 0, 10),
		managers:         make([]ShutdownManager, 0, 3),
		shutdownFinished: make(chan struct{}),
	}
}

func (g *GracefulShutdown) Done() <-chan struct{} {
	return g.shutdownFinished
}

func (g *GracefulShutdown) AddShutdownManager(manager ShutdownManager) {
	g.managers = append(g.managers, manager)
}

func (g *GracefulShutdown) AddShutdownCallback(callback ShutdownCallback) {
	g.callbacks = append(g.callbacks, callback)
}

func (g *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	g.errorHandler = errorHandler
}

func (g *GracefulShutdown) Start() error {
	for _, m := range g.managers {
		if err := m.Start(g); err != nil {
			return err
		}
	}
	return nil
}

func (g *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	g.ReportError(sm.ShutdownStart())
	var wg sync.WaitGroup

	for _, c := range g.callbacks {
		wg.Add(1)
		go func(c ShutdownCallback) {
			defer wg.Done()
			g.ReportError(c.OnShutdown(sm.GetName()))
		}(c)
	}

	wg.Wait()

	g.ReportError(sm.ShutdownFinish())
	close(g.shutdownFinished)
}

func (g *GracefulShutdown) ReportError(err error) {
	if err != nil && g.errorHandler != nil {
		g.errorHandler.OnError(err)
	}
}
