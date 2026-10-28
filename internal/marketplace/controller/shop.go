package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (h *MarketplaceController) context(c *gin.Context) {
	data, err := h.service.Context(c.Request.Context(), actor(c))
	respond(c, data, err)
}

func (h *MarketplaceController) shops(c *gin.Context) {
	data, err := h.service.Shops(c.Request.Context(), actor(c), pagination(c))
	respond(c, data, err)
}

func (h *MarketplaceController) applyShop(c *gin.Context) {
	var in model.ShopInput
	if !bind(c, &in) {
		return
	}
	data, err := h.service.ApplyShop(c.Request.Context(), actor(c), in)
	respond(c, data, err)
}

func (h *MarketplaceController) reviewShop(c *gin.Context) {
	var in model.ShopReview
	if !bind(c, &in) {
		return
	}
	respond(c, nil, h.service.ReviewShop(c.Request.Context(), actor(c), c.Param("id"), in))
}

func (h *MarketplaceController) categories(c *gin.Context) {
	data, err := h.service.Categories(c.Request.Context(), pagination(c))
	respond(c, data, err)
}

func (h *MarketplaceController) category(c *gin.Context) {
	var in struct {
		Name string `json:"name"`
	}
	if !bind(c, &in) {
		return
	}
	data, err := h.service.CreateCategory(c.Request.Context(), actor(c), in.Name)
	respond(c, data, err)
}
