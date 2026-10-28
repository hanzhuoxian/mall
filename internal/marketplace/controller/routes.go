package controller

import "github.com/gin-gonic/gin"

// Register 注册公开及认证后的商城路由。
func (h *MarketplaceController) Register(r *gin.Engine, auth gin.HandlerFunc) {
	public := r.Group("/v1/market")
	public.GET("/catalog", h.catalog)
	public.GET("/catalog/:id", h.catalogProduct)
	public.GET("/categories", h.categories)
	secured := public.Group("", auth, h.authenticated)
	secured.GET("/context", h.context)
	secured.GET("/shops", h.shops)
	secured.POST("/shops", h.applyShop)
	secured.POST("/shops/:id/review", requireAdmin, h.reviewShop)
	secured.POST("/categories", requireAdmin, h.category)
	secured.GET("/products", h.products)
	secured.POST("/products", h.saveProduct)
	secured.PUT("/products/:id", h.saveProduct)
	secured.POST("/products/:id/review", h.reviewProduct)
	secured.POST("/products/:id/stock", h.stock)
	secured.GET("/addresses", h.addresses)
	secured.POST("/addresses", h.address)
	secured.DELETE("/addresses/:id", h.deleteAddress)
	secured.GET("/cart", h.cart)
	secured.PUT("/cart", h.putCart)
	secured.DELETE("/cart/:id", h.deleteCart)
	secured.POST("/checkout", h.checkout)
	secured.GET("/orders", h.orders)
	secured.POST("/orders/:id/actions", h.orderAction)
	secured.GET("/dashboard", h.dashboard)
	secured.GET("/audit", requireAdmin, h.audits)
}
