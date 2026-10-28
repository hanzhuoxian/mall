package marketplace_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hanzhuoxian/mall/internal/marketplace/controller"
	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/service"
	marketmysql "github.com/hanzhuoxian/mall/internal/marketplace/store/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLite covers persistence/rollback behavior. Set MALL_TEST_MYSQL_DSN to an
// isolated test database to run the same assertions against InnoDB row locks.
func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	dialector := sqlite.Open(filepath.Join(t.TempDir(), "market.db"))
	if dsn := os.Getenv("MALL_TEST_MYSQL_DSN"); dsn != "" {
		dialector = mysql.Open(dsn)
	}
	db, err := gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(model.Models()...); err != nil {
		t.Fatal(err)
	}
	sqlDB, _ := db.DB()
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

type fixture struct {
	s        *service.Service
	db       *gorm.DB
	shops    []model.Shop
	products []model.Product
	cart     []model.CartItem
	address  model.Address
	buyer    model.Actor
}

func mustCreate(t *testing.T, db *gorm.DB, value any) {
	t.Helper()
	if err := db.Create(value).Error; err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) fixture {
	t.Helper()
	db := testDB(t)
	f := fixture{s: service.NewService(marketmysql.NewStore(db), true), db: db, buyer: model.Actor{ID: base().ID}}
	cat := model.Category{Base: base(), Name: base().ID}
	if err := db.Create(&cat).Error; err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		shop := model.Shop{
			Base:          base(),
			OwnerID:       base().ID,
			Name:          fmt.Sprintf("shop %d", i),
			Status:        "approved",
			CommissionBPS: 500,
		}
		mustCreate(t, db, &shop)
		p := model.Product{
			Base:       base(),
			ShopID:     shop.ID,
			CategoryID: cat.ID,
			Name:       "Test",
			Status:     "approved",
			Price:      int64(1000 + i*500),
			Stock:      10,
		}
		mustCreate(t, db, &p)
		cart := model.CartItem{Base: base(), UserID: f.buyer.ID, ProductID: p.ID, Quantity: 2}
		mustCreate(t, db, &cart)
		f.shops = append(f.shops, shop)
		f.products = append(f.products, p)
		f.cart = append(f.cart, cart)
	}
	f.address = model.Address{
		Base:      base(),
		UserID:    f.buyer.ID,
		Recipient: "Buyer",
		Phone:     "123456789",
		Detail:    "Test address",
	}
	mustCreate(t, db, &f.address)
	return f
}
func (f fixture) input() model.CheckoutInput {
	return model.CheckoutInput{
		RequestKey: "request",
		AddressID:  f.address.ID,
		CartIDs:    []string{f.cart[0].ID, f.cart[1].ID},
	}
}
func stockOf(t *testing.T, db *gorm.DB, id string) int64 {
	t.Helper()
	var p model.Product
	if err := db.First(&p, "id = ?", id).Error; err != nil {
		t.Fatal(err)
	}
	return p.Stock
}
func TestCheckoutSplitReplayAndLifecycle(t *testing.T) {
	f := setup(t)
	orders, err := f.s.Checkout(context.Background(), f.buyer, f.input())
	if err != nil {
		t.Fatal(err)
	}
	if len(orders) != 2 {
		t.Fatalf("want 2 shops, got %d", len(orders))
	}
	var sum int64
	for _, o := range orders {
		sum += o.Amount
		if o.Commission != o.Amount*500/10000 || len(o.Items) != 1 || o.Address != "Test address" {
			t.Fatalf("wrong snapshot: %+v", o)
		}
	}
	if sum != 5000 {
		t.Fatalf("wrong server price %d", sum)
	}
	for _, p := range f.products {
		if stockOf(t, f.db, p.ID) != 8 {
			t.Fatal("stock not reserved")
		}
	}
	again, err := f.s.Checkout(context.Background(), f.buyer, f.input())
	if err != nil || len(again) != 2 {
		t.Fatalf("retry: %v", err)
	}
	ids := map[string]bool{}
	for _, o := range orders {
		ids[o.ID] = true
	}
	for _, o := range again {
		if !ids[o.ID] {
			t.Fatal("retry created duplicate")
		}
	}
	changed := f.input()
	changed.AddressID = "another"
	if _, err = f.s.Checkout(context.Background(), f.buyer, changed); err == nil {
		t.Fatal("changed idempotent request accepted")
	}
	o := orders[0]
	var shop model.Shop
	f.db.First(&shop, "id = ?", o.ShopID)
	if err = f.s.ActOrder(
		context.Background(),
		model.Actor{ID: "stranger"},
		o.ID,
		model.OrderAction{Action: "cancel"},
	); err != service.ErrForbidden {
		t.Fatalf("cross-tenant action: %v", err)
	}
	if err = f.s.ActOrder(context.Background(), f.buyer, o.ID, model.OrderAction{Action: "pay"}); err != nil {
		t.Fatal(err)
	}
	if err = f.s.ActOrder(
		context.Background(),
		model.Actor{ID: shop.OwnerID},
		o.ID,
		model.OrderAction{Action: "ship", Carrier: "test", Tracking: "123"},
	); err != nil {
		t.Fatal(err)
	}
	if err = f.s.ActOrder(context.Background(), f.buyer, o.ID, model.OrderAction{Action: "receive"}); err != nil {
		t.Fatal(err)
	}
	if err = f.s.ActOrder(
		context.Background(),
		f.buyer,
		o.ID,
		model.OrderAction{Action: "refund", Reason: "wrong item"},
	); err != nil {
		t.Fatal(err)
	}
	if err = f.s.ActOrder(
		context.Background(),
		model.Actor{ID: shop.OwnerID},
		o.ID,
		model.OrderAction{Action: "reject-refund", Reason: "proof required"},
	); err != nil {
		t.Fatal(err)
	}
	var saved model.Order
	f.db.First(&saved, "id = ?", o.ID)
	if saved.Status != "completed" {
		t.Fatal("refund rejection did not restore status")
	}
	other := orders[1]
	if err = f.s.ActOrder(context.Background(), f.buyer, other.ID, model.OrderAction{Action: "cancel"}); err != nil {
		t.Fatal(err)
	}
	if stockOf(t, f.db, other.Items[0].ProductID) != 10 {
		t.Fatal("stock not returned")
	}
	if err = f.s.ActOrder(context.Background(), f.buyer, other.ID, model.OrderAction{Action: "cancel"}); err == nil {
		t.Fatal("double cancel accepted")
	}
	if stockOf(t, f.db, other.Items[0].ProductID) != 10 {
		t.Fatal("stock returned twice")
	}
}
func TestCheckoutRollsBackAllShops(t *testing.T) {
	f := setup(t)
	f.db.Model(&model.Product{}).Where("id = ?", f.products[1].ID).Update("stock", 0)
	if _, err := f.s.Checkout(context.Background(), f.buyer, f.input()); err == nil {
		t.Fatal("oversell accepted")
	}
	if stockOf(t, f.db, f.products[0].ID) != 10 {
		t.Fatal("partial stock commit")
	}
	var count int64
	f.db.Model(&model.Order{}).Where("user_id = ?", f.buyer.ID).Count(&count)
	if count != 0 {
		t.Fatal("partial order commit")
	}
	f.db.Model(&model.CartItem{}).Where("user_id = ?", f.buyer.ID).Count(&count)
	if count != 2 {
		t.Fatal("cart lost on failed checkout")
	}
	f.db.Model(&model.Product{}).Where("id = ?", f.products[1].ID).Update("stock", 10)
	if _, err := f.s.Checkout(context.Background(), f.buyer, f.input()); err != nil {
		t.Fatalf("failed request key cannot be retried: %v", err)
	}
}
func TestExpirationAndRefundReleaseOnce(t *testing.T) {
	f := setup(t)
	orders, err := f.s.Checkout(context.Background(), f.buyer, f.input())
	if err != nil {
		t.Fatal(err)
	}
	f.db.Model(&model.Order{}).Where("id = ?", orders[0].ID).Update("expires_at", time.Now().Add(-time.Hour))
	for i := 0; i < 2; i++ {
		if err = f.s.ExpireOrders(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if stockOf(t, f.db, orders[0].Items[0].ProductID) != 10 {
		t.Fatal("expiration stock mismatch")
	}
	o := orders[1]
	var shop model.Shop
	f.db.First(&shop, "id = ?", o.ShopID)
	for _, action := range []model.OrderAction{{Action: "pay"}, {Action: "refund", Reason: "changed mind"}} {
		if err = f.s.ActOrder(context.Background(), f.buyer, o.ID, action); err != nil {
			t.Fatal(err)
		}
	}
	if err = f.s.ActOrder(
		context.Background(),
		model.Actor{ID: shop.OwnerID},
		o.ID,
		model.OrderAction{Action: "approve-refund"},
	); err != nil {
		t.Fatal(err)
	}
	if stockOf(t, f.db, o.Items[0].ProductID) != 10 {
		t.Fatal("unshipped refund stock mismatch")
	}
	if err = f.s.ActOrder(
		context.Background(),
		model.Actor{ID: shop.OwnerID},
		o.ID,
		model.OrderAction{Action: "approve-refund"},
	); err == nil {
		t.Fatal("duplicate refund accepted")
	}
}
func TestTenantIsolationAndReapproval(t *testing.T) {
	f := setup(t)
	p := f.products[0]
	p.Name = "changed"
	if _, err := f.s.SaveProduct(
		context.Background(),
		model.Actor{ID: f.shops[1].OwnerID},
		p,
	); err != service.ErrForbidden {
		t.Fatalf("foreign merchant edited product: %v", err)
	}
	p.Stock = 999
	out, err := f.s.SaveProduct(context.Background(), model.Actor{ID: f.shops[0].OwnerID}, p)
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "pending" || out.Stock != 10 {
		t.Fatal("edit bypassed review or inventory")
	}
	if _, err = f.s.Checkout(context.Background(), model.Actor{ID: "stranger"}, f.input()); err == nil {
		t.Fatal("foreign cart/address used")
	}
}
func TestHTTPMerchantScopeAndDashboard(t *testing.T) {
	f := setup(t)
	if _, err := f.s.Checkout(context.Background(), f.buyer, f.input()); err != nil {
		t.Fatal(err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := controller.NewMarketplaceController(
		f.s,
		func(_ context.Context, id string) (model.Actor, error) { return model.Actor{ID: id}, nil },
		"user",
	)
	h.Register(r, func(c *gin.Context) { c.Set("user", f.shops[0].OwnerID); c.Next() })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/v1/market/orders?scope=merchant", nil))
	var result struct {
		Data struct {
			Items []model.Order
			Total int64
		}
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || result.Data.Total != 1 || result.Data.Items[0].ShopID != f.shops[0].ID {
		t.Fatalf("scope leaked: %s", w.Body)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/v1/market/dashboard?scope=merchant", nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"orders":1`) {
		t.Fatalf("dashboard: %s", w.Body)
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/v1/market/orders?scope=platform", nil))
	if !strings.Contains(w.Body.String(), `"total":0`) {
		t.Fatal("non-admin saw platform data")
	}
}
func TestMySQLConcurrentReservation(t *testing.T) {
	if os.Getenv("MALL_TEST_MYSQL_DSN") == "" {
		t.Skip("requires isolated MySQL database to verify row locks")
	}
	f := setup(t)
	p := f.products[0]
	f.db.Model(&p).Update("stock", 1)
	type attempt struct {
		in model.CheckoutInput
		a  model.Actor
	}
	attempts := []attempt{}
	for i := 0; i < 2; i++ {
		a := model.Actor{ID: base().ID}
		address := model.Address{Base: base(), UserID: a.ID, Recipient: "x", Phone: "123", Detail: "x"}
		mustCreate(t, f.db, &address)
		cart := model.CartItem{Base: base(), UserID: a.ID, ProductID: p.ID, Quantity: 1}
		mustCreate(t, f.db, &cart)
		attempts = append(
			attempts,
			attempt{
				a:  a,
				in: model.CheckoutInput{AddressID: address.ID, CartIDs: []string{cart.ID}, RequestKey: "concurrent"},
			},
		)
	}
	var wg sync.WaitGroup
	results := make(chan error, 2)
	start := make(chan struct{})
	for _, attempt := range attempts {
		wg.Go(func() { <-start; _, err := f.s.Checkout(context.Background(), attempt.a, attempt.in); results <- err })
	}
	close(start)
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		} else if err != service.ErrStock {
			t.Fatalf("expected stock rejection, got %v", err)
		}
	}
	if success != 1 || stockOf(t, f.db, p.ID) != 0 {
		t.Fatal("concurrent inventory invariant failed")
	}
}

func TestMySQLConcurrentCheckoutReplay(t *testing.T) {
	if os.Getenv("MALL_TEST_MYSQL_DSN") == "" {
		t.Skip("requires isolated MySQL database")
	}
	f := setup(t)
	const clients = 8
	start := make(chan struct{})
	type result struct {
		err    error
		orders []model.Order
	}
	results := make(chan result, clients)
	var wg sync.WaitGroup
	for range clients {
		wg.Go(func() {
			<-start
			orders, err := f.s.Checkout(context.Background(), f.buyer, f.input())
			results <- result{orders: orders, err: err}
		})
	}
	close(start)
	wg.Wait()
	close(results)
	var batch string
	for r := range results {
		if r.err != nil {
			t.Fatalf("concurrent replay failed: %v", r.err)
		}
		if len(r.orders) != 2 {
			t.Fatalf("expected 2 orders, got %d", len(r.orders))
		}
		if batch == "" {
			batch = r.orders[0].CheckoutID
		}
		for _, o := range r.orders {
			if o.CheckoutID != batch {
				t.Fatal("duplicate checkout created")
			}
		}
	}
	var count int64
	if err := f.db.Model(&model.Order{}).Where("user_id = ?", f.buyer.ID).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("duplicate orders: %d", count)
	}
	for _, p := range f.products {
		if stockOf(t, f.db, p.ID) != 8 {
			t.Fatal("replayed request reserved stock more than once")
		}
	}
}

func TestMySQLCancelRacesExpiration(t *testing.T) {
	if os.Getenv("MALL_TEST_MYSQL_DSN") == "" {
		t.Skip("requires isolated MySQL database")
	}
	f := setup(t)
	orders, err := f.s.Checkout(context.Background(), f.buyer, f.input())
	if err != nil {
		t.Fatal(err)
	}
	o := orders[0]
	if err = f.db.Model(&model.Order{}).
		Where("id = ?", o.ID).
		Update("expires_at", time.Now().Add(-time.Minute)).
		Error; err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	results := make(chan error, 2)
	wg.Go(func() {
		<-start
		results <- f.s.ActOrder(context.Background(), f.buyer, o.ID, model.OrderAction{Action: "cancel"})
	})
	wg.Go(func() { <-start; results <- f.s.ExpireOrders(context.Background()) })
	close(start)
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil && err != service.ErrInvalid {
			t.Fatalf("unexpected racing error: %v", err)
		}
	}
	var saved model.Order
	if err = f.db.First(&saved, "id = ?", o.ID).Error; err != nil {
		t.Fatal(err)
	}
	if saved.Status != model.OrderCanceled || stockOf(t, f.db, o.Items[0].ProductID) != 10 {
		t.Fatal("racing cancellation returned stock incorrectly")
	}
}

func base() model.Base { return model.Base{ID: uuid.NewString()} }

func TestCatalogByIDVisibility(t *testing.T) {
	f := setup(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := controller.NewMarketplaceController(f.s, func(_ context.Context, id string) (model.Actor, error) {
		return model.Actor{ID: id, Admin: id == "platform-admin"}, nil
	}, "user")
	h.Register(router, func(c *gin.Context) { c.Set("user", c.GetHeader("X-Test-User")); c.Next() })
	call := func(id string, want int) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest("GET", "/v1/market/catalog/"+id, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != want {
			t.Fatalf("GET catalog/%s: HTTP %d want %d: %s", id, w.Code, want, w.Body)
		}
		return w
	}
	visible := call(f.products[0].ID, 200)
	var body struct {
		Data model.Product `json:"data"`
	}
	if err := json.Unmarshal(visible.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Data.ID != f.products[0].ID {
		t.Fatalf("wrong product returned: %+v", body.Data)
	}

	pending := model.Product{Base: base(), ShopID: f.shops[0].ID, CategoryID: f.products[0].CategoryID, Name: "Pending", Status: "pending"}
	mustCreate(t, f.db, &pending)
	call(pending.ID, 404)

	offlineShop := model.Shop{Base: base(), OwnerID: base().ID, Name: "Offline", Status: "pending"}
	mustCreate(t, f.db, &offlineShop)
	offlineProduct := model.Product{Base: base(), ShopID: offlineShop.ID, CategoryID: f.products[0].CategoryID, Name: "Hidden", Status: "approved"}
	mustCreate(t, f.db, &offlineProduct)
	call(offlineProduct.ID, 404)

	if err := f.db.Model(&model.Product{}).Where("id = ?", f.products[1].ID).Update("status", "offline").Error; err != nil {
		t.Fatal(err)
	}
	call(f.products[1].ID, 404)
}

func TestHTTPCatalogAndCheckoutWorkflow(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	h := controller.NewMarketplaceController(f.s, func(_ context.Context, id string) (model.Actor, error) {
		return model.Actor{ID: id, Admin: id == "platform-admin"}, nil
	}, "user")
	h.Register(router, func(c *gin.Context) { c.Set("user", c.GetHeader("X-Test-User")); c.Next() })
	call := func(user, method, path string, input any, want int, output any) {
		t.Helper()
		data, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(method, "/v1/market"+path, strings.NewReader(string(data)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User", user)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var body struct {
			Data json.RawMessage
			Code int
		}
		if err = json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if w.Code != want {
			t.Fatalf("%s %s: HTTP %d want %d: %s", method, path, w.Code, want, w.Body)
		}
		if want == 200 && body.Code != 0 {
			t.Fatalf("success code: %s", w.Body)
		}
		if want != 200 && body.Code != want {
			t.Fatalf("error code: %s", w.Body)
		}
		if output != nil {
			if err = json.Unmarshal(body.Data, output); err != nil {
				t.Fatal(err)
			}
		}
	}
	const merchant = "new-merchant"
	var shop model.Shop
	call(
		merchant,
		"POST",
		"/shops",
		model.ShopInput{Name: "New shop", Contact: "contact", License: "license"},
		200,
		&shop,
	)
	call(
		f.buyer.ID,
		"POST",
		"/shops/"+shop.ID+"/review",
		model.ShopReview{Status: "approved", CommissionBPS: 500},
		403,
		nil,
	)
	call(
		"platform-admin",
		"POST",
		"/shops/"+shop.ID+"/review",
		model.ShopReview{Status: "approved", CommissionBPS: 500},
		200,
		nil,
	)
	var category model.Category
	call("platform-admin", "POST", "/categories", map[string]string{"name": "HTTP category"}, 200, &category)
	var product model.Product
	call(
		merchant,
		"POST",
		"/products",
		model.Product{ShopID: shop.ID, CategoryID: category.ID, Name: "HTTP product", Price: 2500, Stock: 3},
		200,
		&product,
	)
	var catalog model.Page[model.Product]
	call("", "GET", "/catalog?categoryId="+category.ID, nil, 200, &catalog)
	if catalog.Total != 0 {
		t.Fatal("pending product is publicly visible")
	}
	call("platform-admin", "POST", "/products/"+product.ID+"/review", model.ProductReview{Status: "approved"}, 200, nil)
	call("", "GET", "/catalog?categoryId="+category.ID, nil, 200, &catalog)
	if catalog.Total != 1 || catalog.Items[0].ID != product.ID {
		t.Fatal("approved catalog mismatch")
	}
	call("stranger", "POST", "/products/"+product.ID+"/stock", map[string]int{"delta": 1}, 403, nil)
	call(merchant, "POST", "/products/"+product.ID+"/stock", map[string]int{"delta": 1}, 200, nil)
	call(f.buyer.ID, "PUT", "/cart", model.CartItem{ProductID: product.ID, Quantity: 2}, 200, nil)
	var cart model.Page[model.CartItem]
	call(f.buyer.ID, "GET", "/cart", nil, 200, &cart)
	var cartID string
	for _, item := range cart.Items {
		if item.ProductID == product.ID {
			cartID = item.ID
		}
	}
	if cartID == "" {
		t.Fatal("cart update missing")
	}
	var address model.Address
	call(
		f.buyer.ID,
		"POST",
		"/addresses",
		model.Address{Recipient: "Buyer", Phone: "123", Detail: "Address"},
		200,
		&address,
	)
	var orders []model.Order
	input := model.CheckoutInput{AddressID: address.ID, CartIDs: []string{cartID}, RequestKey: "http-checkout"}
	call(f.buyer.ID, "POST", "/checkout", input, 200, &orders)
	if len(orders) != 1 || orders[0].Amount != 5000 || stockOf(t, f.db, product.ID) != 2 {
		t.Fatalf("incorrect checkout: %+v", orders)
	}
	call(f.buyer.ID, "POST", "/orders/"+orders[0].ID+"/actions", model.OrderAction{Action: "pay"}, 200, nil)
	call(
		merchant,
		"POST",
		"/orders/"+orders[0].ID+"/actions",
		model.OrderAction{Action: "ship", Carrier: "TEST", Tracking: "TEST123"},
		200,
		nil,
	)
	call(f.buyer.ID, "POST", "/orders/"+orders[0].ID+"/actions", model.OrderAction{Action: "receive"}, 200, nil)
	var summary model.Summary
	call(merchant, "GET", "/dashboard?scope=merchant", nil, 200, &summary)
	if summary.GMV != 5000 || summary.Commission != 250 || summary.Payable != 4750 {
		t.Fatalf("incorrect summary: %+v", summary)
	}
	// Service permission guards also apply when invoked without the HTTP adapter.
	if err := f.s.ReviewShop(ctx, f.buyer, shop.ID, model.ShopReview{Status: "approved"}); err != service.ErrForbidden {
		t.Fatalf("service permission: %v", err)
	}
	call(f.buyer.ID, "GET", "/audit", nil, 403, nil)
	call("platform-admin", "GET", "/audit", nil, 200, nil)
}
