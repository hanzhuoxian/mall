package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (h *MarketplaceController) checkout(c *gin.Context) {
	var in model.CheckoutInput
	if !bind(c, &in) {
		return
	}
	data, err := h.service.Checkout(c.Request.Context(), actor(c), in)
	respond(c, data, err)
}

func (h *MarketplaceController) orders(c *gin.Context) {
	data, err := h.service.Orders(c.Request.Context(), actor(c), c.Query("scope"), c.Query("status"), pagination(c))
	respond(c, data, err)
}

func (h *MarketplaceController) orderAction(c *gin.Context) {
	var in model.OrderAction
	if !bind(c, &in) {
		return
	}
	respond(c, nil, h.service.ActOrder(c.Request.Context(), actor(c), c.Param("id"), in))
}

func (h *MarketplaceController) dashboard(c *gin.Context) {
	data, err := h.service.Dashboard(c.Request.Context(), actor(c), c.Query("scope"))
	respond(c, data, err)
}

func (h *MarketplaceController) audits(c *gin.Context) {
	data, err := h.service.Audits(c.Request.Context(), actor(c), pagination(c))
	respond(c, data, err)
}
