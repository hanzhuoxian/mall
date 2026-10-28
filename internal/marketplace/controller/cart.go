package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (h *MarketplaceController) addresses(c *gin.Context) {
	data, err := h.service.Addresses(c.Request.Context(), actor(c), pagination(c))
	respond(c, data, err)
}

func (h *MarketplaceController) address(c *gin.Context) {
	var in model.Address
	if !bind(c, &in) {
		return
	}
	data, err := h.service.CreateAddress(c.Request.Context(), actor(c), in)
	respond(c, data, err)
}

func (h *MarketplaceController) deleteAddress(c *gin.Context) {
	respond(c, nil, h.service.DeleteAddress(c.Request.Context(), actor(c), c.Param("id")))
}

func (h *MarketplaceController) cart(c *gin.Context) {
	data, err := h.service.Cart(c.Request.Context(), actor(c), pagination(c))
	respond(c, data, err)
}

func (h *MarketplaceController) putCart(c *gin.Context) {
	var in model.CartItem
	if !bind(c, &in) {
		return
	}
	respond(c, nil, h.service.PutCart(c.Request.Context(), actor(c), in))
}

func (h *MarketplaceController) deleteCart(c *gin.Context) {
	respond(c, nil, h.service.DeleteCart(c.Request.Context(), actor(c), c.Param("id")))
}
