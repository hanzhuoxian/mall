package userserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzhuoxian/mall/internal/userserver/cache"
)

func initRouter(g *gin.Engine) {
	RegisterMiddleware(g)
	RegisterController(g)
}

func RegisterMiddleware(g *gin.Engine) {}
func RegisterController(g *gin.Engine) {
	g.GET("/hello", func(ctx *gin.Context) {
		u := cache.Get().User()
		u.SetUser(ctx, "username", "hanjian")
		username, err := u.GetUser(ctx, "username")
		if err != nil {
			ctx.JSON(http.StatusNotFound, "")
		} else {
			u.GetUser(ctx, username)
		}
	})
}
