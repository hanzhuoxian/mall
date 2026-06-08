package userserver

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/userserver/controller"
)

// initRouter 完成路由的整体初始化，依次注册中间件和控制器路由。
func initRouter(g *gin.Engine, ctrls *controller.Controllers) {
	registerMiddleware(g)
	registerController(g, ctrls)
}

// registerMiddleware 注册全局中间件，当前为空占位，后续可在此扩展。
func registerMiddleware(_ *gin.Engine) {}

// registerController 将各控制器的处理函数注册到对应的 HTTP 路由上。
func registerController(g *gin.Engine, ctrls *controller.Controllers) {
	g.GET("/hello", ctrls.User.Hello)
}
