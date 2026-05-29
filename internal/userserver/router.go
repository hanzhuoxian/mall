package userserver

import (
	"github.com/gin-gonic/gin"
	"github.com/hanzhuoxian/mall/internal/userserver/controller"
)

func initRouter(g *gin.Engine, ctrls *controller.Controllers) {
	registerMiddleware(g)
	registerController(g, ctrls)
}

func registerMiddleware(_ *gin.Engine) {}

func registerController(g *gin.Engine, ctrls *controller.Controllers) {
	g.GET("/hello", ctrls.User.Hello)
}
