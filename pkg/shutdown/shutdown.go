// Package shutdown 提供优雅关闭（Graceful Shutdown）机制，支持注册多个关闭管理器和回调函数，
// 在服务关闭时并发执行所有回调并等待完成。
package shutdown

import (
	"sync"
)

// ShutdownCallback 表示一个关闭回调，在服务关闭时被调用。
// 参数为触发关闭的管理器名称，返回执行过程中的错误。
type ShutdownCallback interface {
	OnShutdown(string) error
}

// ShutdownFunc 是 ShutdownCallback 的函数类型适配器，
// 允许将普通函数直接作为 ShutdownCallback 使用。
type ShutdownFunc func(string) error

// OnShutdown 实现 ShutdownCallback 接口。
func (f ShutdownFunc) OnShutdown(s string) error {
	return f(s)
}

// ErrorHandler 用于处理关闭过程中产生的错误。
type ErrorHandler interface {
	OnError(err error)
}

// ErrorFunc 是 ErrorHandler 的函数类型适配器。
type ErrorFunc func(err error)

// OnError 实现 ErrorHandler 接口。
func (f ErrorFunc) OnError(err error) {
	f(err)
}

// ShutdownManager 负责监听关闭信号并驱动整个关闭流程。
// 典型实现包括：监听 OS 信号、HTTP 服务关闭等。
type ShutdownManager interface {
	// GetName 返回管理器名称，作为回调参数传入，用于区分触发来源。
	GetName() string
	// Start 启动管理器，开始监听关闭信号。
	Start(gs GracefulShutter) error
	// ShutdownStart 在执行回调前调用，可用于准备工作（如停止接收新请求）。
	ShutdownStart() error
	// ShutdownFinish 在所有回调执行完毕后调用，用于最终清理。
	ShutdownFinish() error
}

// GracefulShutter 是 GracefulShutdown 对外暴露给 ShutdownManager 的接口，
// 避免管理器持有完整的 GracefulShutdown 引用。
type GracefulShutter interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error)
	AddShutdownCallback(sc ShutdownCallback)
}

// GracefulShutdown 是优雅关闭的核心结构，管理回调列表、关闭管理器和错误处理器。
type GracefulShutdown struct {
	callbacks        []ShutdownCallback
	managers         []ShutdownManager
	errorHandler     ErrorHandler
	shutdownFinished chan struct{} // 所有回调执行完毕后关闭，用于通知外部等待方
}

// New 创建并返回一个新的 GracefulShutdown 实例。
func New() *GracefulShutdown {
	return &GracefulShutdown{
		callbacks:        make([]ShutdownCallback, 0, 10),
		managers:         make([]ShutdownManager, 0, 3),
		shutdownFinished: make(chan struct{}),
	}
}

// Done 返回一个 channel，当所有关闭回调执行完毕后该 channel 会被关闭。
// 可用于阻塞主 goroutine 直到关闭流程结束。
func (g *GracefulShutdown) Done() <-chan struct{} {
	return g.shutdownFinished
}

// AddShutdownManager 注册一个关闭管理器，Start 时会逐一启动。
func (g *GracefulShutdown) AddShutdownManager(manager ShutdownManager) {
	g.managers = append(g.managers, manager)
}

// AddShutdownCallback 注册一个关闭回调，在关闭时并发执行。
func (g *GracefulShutdown) AddShutdownCallback(callback ShutdownCallback) {
	g.callbacks = append(g.callbacks, callback)
}

// SetErrorHandler 设置错误处理器，用于处理关闭过程中的非致命错误。
func (g *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	g.errorHandler = errorHandler
}

// Start 启动所有已注册的关闭管理器，使其开始监听关闭信号。
func (g *GracefulShutdown) Start() error {
	for _, m := range g.managers {
		if err := m.Start(g); err != nil {
			return err
		}
	}
	return nil
}

// StartShutdown 由 ShutdownManager 在收到关闭信号时调用，并发执行所有回调，
// 等待全部完成后通知 Done() channel。
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

// ReportError 将错误传递给已注册的错误处理器，若未设置处理器则忽略。
func (g *GracefulShutdown) ReportError(err error) {
	if err != nil && g.errorHandler != nil {
		g.errorHandler.OnError(err)
	}
}
