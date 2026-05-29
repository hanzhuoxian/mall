package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

const Name = "PosixSignalManager"

type PosixSignalManager struct {
	signals []os.Signal
}

func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = []os.Signal{os.Interrupt, syscall.SIGTERM}
	}

	return &PosixSignalManager{signals: sig}
}

func (p *PosixSignalManager) GetName() string {
	return Name
}

func (p *PosixSignalManager) Start(gs GracefulShutter) error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, p.signals...)
		<-c

		gs.StartShutdown(p)
	}()
	return nil
}

func (p *PosixSignalManager) ShutdownStart() error {
	return nil
}

func (p *PosixSignalManager) ShutdownFinish() error {
	return nil
}
