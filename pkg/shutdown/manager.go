package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Name 是 PosixSignalManager 的唯一标识名称，作为回调参数传入，用于区分关闭来源。
const Name = "PosixSignalManager"

// PosixSignalManager 是基于 POSIX 信号的关闭管理器，监听操作系统信号（如 SIGINT、SIGTERM）
// 并触发优雅关闭流程。
type PosixSignalManager struct {
	signals []os.Signal
}

// NewPosixSignalManager 创建一个 PosixSignalManager 实例。
// 若未指定信号，默认监听 os.Interrupt（SIGINT）和 SIGTERM。
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = []os.Signal{os.Interrupt, syscall.SIGTERM}
	}

	return &PosixSignalManager{signals: sig}
}

// GetName 返回管理器名称，实现 ShutdownManager 接口。
func (p *PosixSignalManager) GetName() string {
	return Name
}

// Start 在独立 goroutine 中监听 OS 信号，收到信号后触发 GracefulShutter.StartShutdown。
func (p *PosixSignalManager) Start(gs GracefulShutter) error {
	ctx, stop := signal.NotifyContext(context.Background(), p.signals...)

	go func() {
		<-ctx.Done()
		stop()

		gs.StartShutdown(p)
	}()
	return nil
}

// ShutdownStart 在回调执行前调用，PosixSignalManager 无需额外准备，直接返回 nil。
func (p *PosixSignalManager) ShutdownStart() error {
	return nil
}

// ShutdownFinish 在所有回调执行完毕后调用，PosixSignalManager 无需清理工作，直接返回 nil。
func (p *PosixSignalManager) ShutdownFinish() error {
	return nil
}
