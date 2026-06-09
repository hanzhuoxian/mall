// Package cache 定义了用户服务缓存层的抽象接口及全局单例管理函数。
package cache

// Factory 是缓存工厂接口，提供各缓存子接口的访问入口及资源关闭能力。
type Factory interface {
	User() UserCache
	Close() error
}

// cache 是全局缓存工厂单例，由 Set 初始化，由 Get 读取。
var cache Factory

// Set 设置全局缓存工厂实例，通常在服务启动时由依赖注入框架调用。
func Set(f Factory) {
	cache = f
}

// Get 返回当前全局缓存工厂实例。
func Get() Factory {
	return cache
}
