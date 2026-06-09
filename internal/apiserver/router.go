package apiserver

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/apiserver/controller"
)

func installRoutes(r *gin.Engine, controllers *controller.Controllers) {
	v1 := r.Group("/v1")
	{
		users := v1.Group("/users")
		users.GET("", controllers.User.ListUsers)
		users.GET("/:id", controllers.User.GetUser)
	}
}
