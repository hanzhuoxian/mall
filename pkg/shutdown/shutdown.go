package shutdown

type ShutdownCallback interface {
	OnShutdown(string) error
}

type ShutdownFunc func(string) error

type ErrorHandler interface {
	OnError(err error)
}

type ErrorFunc func(err error)

func (f ErrorFunc) OnError(err error) {
	f(err)
}

type ShutdownManager interface {
	GetName() string
	Start(gs GracefulShutdowner) error
	ShutdownStart() error
	ShutdownFinish() error
}

type GracefulShutdowner interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(sc ShutdownCallback)
}

type GracefulShotdown struct {
	callbacks    []ShutdownCallback
	managers     []ShutdownManager
	errorHandler ErrorHandler
}
