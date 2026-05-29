package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hanzhuoxian/mall/internal/userserver/cache"
)

type UserController struct {
	cache cache.Factory
}

func NewUserController(c cache.Factory) *UserController {
	return &UserController{cache: c}
}

func (uc *UserController) Hello(ctx *gin.Context) {
	u := uc.cache.User()
	if err := u.SetUser(ctx, "username", "hanjian"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	username, err := u.GetUser(ctx, "username")
	if err != nil {
		ctx.JSON(http.StatusNotFound, "")
		return
	}
	ctx.JSON(http.StatusOK, username)
}
