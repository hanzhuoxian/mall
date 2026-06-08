// Package store 定义了用户服务数据存储层的抽象接口及全局单例管理函数。
package store

// Factory 是数据存储工厂接口，提供资源关闭能力。
// 后续可在此接口中扩展各业务实体的 Store 访问方法。
type Factory interface {
	Close() error
}

// client 是全局存储工厂单例，由 Set 初始化，由 Get 读取。
var client Factory

// Set 设置全局存储工厂实例，通常在服务启动时由依赖注入框架调用。
func Set(f Factory) {
	client = f
}

// Get 返回当前全局存储工厂实例。
func Get() Factory {
	return client
}
