package apiserver

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/apiserver/controller"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware"
	"github.com/hanzhuoxian/mall/internal/pkg/middleware/auth"
)

func installRoutes(r *gin.Engine, controllers *controller.Controllers, authStrategy middleware.AuthStrategy, jwtStrategy auth.JWTStrategy) {
	r.POST("/login", jwtStrategy.LoginHandler)
	r.POST("/logout", jwtStrategy.LogoutHandler)
	r.POST("/refresh", jwtStrategy.RefreshHandler)

	v1 := r.Group("/v1")
	{
		users := v1.Group("/users", authStrategy.AuthFunc())
		users.GET("", controllers.User.ListUsers)
		users.GET("/:id", controllers.User.GetUser)
	}
}
