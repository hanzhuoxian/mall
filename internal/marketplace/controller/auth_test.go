package controller

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/service"
)

func TestProtectedEndpointsRejectAnonymous(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewMarketplaceController(
		service.NewService(nil, false),
		func(context.Context, string) (model.Actor, error) { return model.Actor{}, service.ErrForbidden },
		"user",
	)
	h.Register(r, func(c *gin.Context) { c.Next() })
	for _, route := range []struct{ method, path string }{{"POST", "/shops"}, {"GET", "/orders"}, {"POST", "/checkout"}, {"GET", "/audit"}, {"POST", "/categories"}, {"PUT", "/cart"}, {"POST", "/products/x/stock"}} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(route.method, "/v1/market"+route.path, nil))
		if w.Code != 401 {
			t.Fatalf("%s: %d", route.path, w.Code)
		}
	}
}
func TestCustomerCannotAdminister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewMarketplaceController(
		service.NewService(nil, false),
		func(context.Context, string) (model.Actor, error) { return model.Actor{ID: "buyer"}, nil },
		"user",
	)
	h.Register(r, func(c *gin.Context) { c.Set("user", "buyer"); c.Next() })
	for _, route := range []struct{ method, path string }{{"POST", "/shops/x/review"}, {"POST", "/categories"}, {"GET", "/audit"}} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(route.method, "/v1/market"+route.path, nil))
		if w.Code != 403 {
			t.Fatalf("%s: %d", route.path, w.Code)
		}
	}
}
