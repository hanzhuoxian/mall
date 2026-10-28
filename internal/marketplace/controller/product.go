package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func (h *MarketplaceController) catalog(c *gin.Context) {
	data, err := h.service.Catalog(c.Request.Context(), c.Query("name"), c.Query("categoryId"), pagination(c))
	respond(c, data, err)
}

func (h *MarketplaceController) catalogProduct(c *gin.Context) {
	data, err := h.service.CatalogProduct(c.Request.Context(), c.Param("id"))
	respond(c, data, err)
}

func (h *MarketplaceController) products(c *gin.Context) {
	data, err := h.service.Products(c.Request.Context(), actor(c), c.Query("name"), c.Query("status"), pagination(c))
	respond(c, data, err)
}

func (h *MarketplaceController) saveProduct(c *gin.Context) {
	var in model.Product
	if !bind(c, &in) {
		return
	}
	in.ID = c.Param("id")
	data, err := h.service.SaveProduct(c.Request.Context(), actor(c), in)
	respond(c, data, err)
}

func (h *MarketplaceController) reviewProduct(c *gin.Context) {
	var in model.ProductReview
	if !bind(c, &in) {
		return
	}
	respond(c, nil, h.service.ReviewProduct(c.Request.Context(), actor(c), c.Param("id"), in))
}

func (h *MarketplaceController) stock(c *gin.Context) {
	var in struct {
		Delta int64 `json:"delta"`
	}
	if !bind(c, &in) {
		return
	}
	respond(c, nil, h.service.AdjustStock(c.Request.Context(), actor(c), c.Param("id"), in.Delta))
}
