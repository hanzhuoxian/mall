package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	marketmodel "github.com/hanzhuoxian/mall/internal/marketplace/model"
	usermodel "github.com/hanzhuoxian/mall/internal/userserver/model"
	"github.com/hanzhuoxian/mall/pkg/auth"
)

// This opt-in test runs the actual API and gRPC user binaries against an isolated
// MySQL schema and Redis DB 15. It never resets a Redis database or drops a schema.
func TestMySQLRedisStack(t *testing.T) {
	userBin, apiBin, dsn := os.Getenv(
		"MALL_TEST_USER_BINARY",
	), os.Getenv(
		"MALL_TEST_API_BINARY",
	), os.Getenv(
		"MALL_TEST_MYSQL_DSN",
	)
	if userBin == "" || apiBin == "" || dsn == "" {
		t.Skip("requires built binaries and isolated MySQL/Redis")
	}
	cfg, err := mysqldriver.ParseDSN(dsn)
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, _ := db.DB()
	t.Cleanup(func() { _ = sqlDB.Close() })
	models := append(
		marketmodel.Models(),
		&usermodel.SysUser{},
		&usermodel.SysRole{},
		&usermodel.SysPermission{},
		&usermodel.SysUserRole{},
		&usermodel.SysRolePermission{},
	)
	if err = db.AutoMigrate(models...); err != nil {
		t.Fatal(err)
	}
	const password = "Integration-Test-Only-Password-2026!"
	hash, err := auth.GeneratePassword(password)
	if err != nil {
		t.Fatal(err)
	}
	admin := usermodel.SysUser{
		ObjectMeta: usermodel.ObjectMeta{InstanceID: uuid.NewString(), Name: "admin"},
		Username:   "admin",
		Email:      "stack-admin@example.test",
		Nickname:   "Admin",
		Password:   hash,
		UserStatus: 1,
	}
	if err = db.Where("username = ?", "admin").FirstOrCreate(&admin).Error; err != nil {
		t.Fatal(err)
	}
	freePort := func() string {
		l, listenErr := net.Listen("tcp", "127.0.0.1:0")
		if listenErr != nil {
			t.Fatal(listenErr)
		}
		port := strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
		_ = l.Close()
		return port
	}
	userHTTP, userGRPC, apiHTTP := freePort(), freePort(), freePort()
	redisAddr := os.Getenv("MALL_TEST_REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	common := []string{
		"--mysql.host",
		cfg.Addr,
		"--mysql.username",
		cfg.User,
		"--mysql.password",
		cfg.Passwd,
		"--mysql.database",
		cfg.DBName,
		"--redis.addrs",
		redisAddr,
		"--redis.db",
		"15",
		"--server.healthz",
	}
	start := func(binary, port string, args ...string) {
		log, startErr := os.Create(filepath.Join(t.TempDir(), "server.log"))
		if startErr != nil {
			t.Fatal(startErr)
		}
		cmd := exec.Command(binary, append(append([]string{}, common...), args...)...)
		cmd.Stdout = log
		cmd.Stderr = log
		if startErr = cmd.Start(); startErr != nil {
			t.Fatal(startErr)
		}
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		t.Cleanup(func() {
			_ = cmd.Process.Signal(os.Interrupt)
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				_ = cmd.Process.Kill()
				<-done
			}
			_ = log.Close()
		})
		client := http.Client{Timeout: time.Second}
		deadline := time.Now().Add(25 * time.Second)
		for time.Now().Before(deadline) {
			resp, requestErr := client.Get("http://127.0.0.1:" + port + "/healthz")
			if requestErr == nil {
				resp.Body.Close()
				if resp.StatusCode == 200 {
					return
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
		data, _ := os.ReadFile(log.Name())
		t.Fatalf("server not ready: %s", data)
	}
	start(
		userBin,
		userHTTP,
		"--grpc.bind-address",
		"127.0.0.1",
		"--grpc.bind-port",
		userGRPC,
		"--insecure.bind-port",
		userHTTP,
	)
	start(
		apiBin,
		apiHTTP,
		"--grpc.bind-address",
		"127.0.0.1",
		"--grpc.bind-port",
		userGRPC,
		"--insecure.bind-port",
		apiHTTP,
		"--server.jwt-secret",
		uuid.NewString()+uuid.NewString(),
		"--sandbox-payment",
	)
	client := http.Client{Timeout: 10 * time.Second}
	baseURL := "http://127.0.0.1:" + apiHTTP
	call := func(token, method, path string, input any, want int) json.RawMessage {
		data, _ := json.Marshal(input)
		req, requestErr := http.NewRequest(method, baseURL+path, bytes.NewReader(data))
		if requestErr != nil {
			t.Fatal(requestErr)
		}
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, requestErr := client.Do(req)
		if requestErr != nil {
			t.Fatal(requestErr)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != want {
			t.Fatalf("%s %s: got HTTP %d, want %d: %s", method, path, resp.StatusCode, want, body)
		}
		var envelope struct {
			Data json.RawMessage
			Code int
		}
		if requestErr = json.Unmarshal(body, &envelope); requestErr != nil {
			t.Fatal(requestErr)
		}
		if want == 200 && envelope.Code != 0 {
			t.Fatalf("%s failed: %s", path, body)
		}
		return envelope.Data
	}
	decode := func(raw json.RawMessage, out any) {
		if decodeErr := json.Unmarshal(raw, out); decodeErr != nil {
			t.Fatal(decodeErr)
		}
	}
	login := func(name string) string {
		var token struct {
			AccessToken string `json:"access_token"`
		}
		decode(call("", "POST", "/login", map[string]any{"identifier": name, "password": password}, 200), &token)
		if token.AccessToken == "" {
			t.Fatal("missing token")
		}
		return token.AccessToken
	}
	adminToken := login("admin")
	call("", "GET", "/v1/market/orders", nil, 401)
	call(adminToken, "GET", "/v1/users/me", nil, 200)
	suffix := uuid.NewString()[:8]
	tokens := []string{}
	userIDs := []string{}
	for _, role := range []string{"merchant_a_", "merchant_b_", "buyer_"} {
		name := role + suffix
		var result struct {
			User struct {
				InstanceID string `json:"instanceId"`
			}
		}
		decode(
			call(
				adminToken,
				"POST",
				"/v1/users",
				map[string]any{
					"name":     name,
					"username": name,
					"nickname": name,
					"email":    name + "@example.test",
					"password": password,
				},
				200,
			),
			&result,
		)
		tokens = append(tokens, login(name))
		userIDs = append(userIDs, result.User.InstanceID)
	}

	// Read only this test's generated challenge from its isolated Redis database.
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr, DB: 15})
	t.Cleanup(func() { _ = rdb.Close() })
	var challenge struct {
		ID    string
		Image string
	}
	decode(call("", "GET", "/captcha", nil, 200), &challenge)
	answer, err := rdb.Get(context.Background(), "captcha:"+challenge.ID).Result()
	if err != nil {
		t.Fatal(err)
	}
	registration := map[string]any{
		"username":    "registered_" + suffix,
		"email":       "registered_" + suffix + "@example.test",
		"password":    password,
		"captchaId":   challenge.ID,
		"captchaCode": answer,
	}
	call("", "POST", "/register", registration, 200)
	call("", "POST", "/register", registration, 400)
	registeredToken := login("registered_" + suffix)
	var registered struct {
		User struct {
			InstanceID string `json:"instanceId"`
		}
	}
	decode(call(registeredToken, "GET", "/v1/users/me", nil, 200), &registered)
	tokens[2] = registeredToken
	userIDs[2] = registered.User.InstanceID
	buyer := tokens[2]
	call(buyer, "GET", "/v1/users", nil, 403)
	var category marketmodel.Category
	decode(call(adminToken, "POST", "/v1/market/categories", map[string]any{"name": "stack-" + suffix}, 200), &category)
	shops := []marketmodel.Shop{}
	products := []marketmodel.Product{}
	for i := 0; i < 2; i++ {
		var shop marketmodel.Shop
		decode(
			call(
				tokens[i],
				"POST",
				"/v1/market/shops",
				map[string]any{
					"name":    fmt.Sprintf("Stack shop %d %s", i, suffix),
					"contact": "test contact",
					"license": "TEST-" + suffix,
				},
				200,
			),
			&shop,
		)
		call(
			adminToken,
			"POST",
			"/v1/market/shops/"+shop.ID+"/review",
			map[string]any{"status": "approved", "commissionBps": 500},
			200,
		)
		var product marketmodel.Product
		decode(
			call(
				tokens[i],
				"POST",
				"/v1/market/products",
				map[string]any{
					"shopId":     shop.ID,
					"categoryId": category.ID,
					"name":       "Stack product",
					"price":      1990,
					"stock":      5,
				},
				200,
			),
			&product,
		)
		call(adminToken, "POST", "/v1/market/products/"+product.ID+"/review", map[string]any{"status": "approved"}, 200)
		call(buyer, "PUT", "/v1/market/cart", map[string]any{"productId": product.ID, "quantity": 1}, 200)
		shops = append(shops, shop)
		products = append(products, product)
	}
	call(tokens[1], "POST", "/v1/market/products/"+products[0].ID+"/stock", map[string]any{"delta": 1}, 403)
	var address marketmodel.Address
	decode(
		call(
			buyer,
			"POST",
			"/v1/market/addresses",
			map[string]any{"recipient": "Stack buyer", "phone": "123456789", "detail": "Integration test address"},
			200,
		),
		&address,
	)
	var cart struct{ Items []marketmodel.CartItem }
	decode(call(buyer, "GET", "/v1/market/cart", nil, 200), &cart)
	ids := []string{}
	for _, item := range cart.Items {
		ids = append(ids, item.ID)
	}
	input := map[string]any{"addressId": address.ID, "cartIds": ids, "requestKey": "stack-" + suffix}
	var orders, replayed []marketmodel.Order
	decode(call(buyer, "POST", "/v1/market/checkout", input, 200), &orders)
	decode(call(buyer, "POST", "/v1/market/checkout", input, 200), &replayed)
	if len(orders) != 2 || len(replayed) != 2 || orders[0].CheckoutID != replayed[0].CheckoutID {
		t.Fatal("split or idempotency failed")
	}
	for _, o := range orders {
		merchant := tokens[0]
		if o.ShopID == shops[1].ID {
			merchant = tokens[1]
		}
		call(buyer, "POST", "/v1/market/orders/"+o.ID+"/actions", map[string]any{"action": "pay"}, 200)
		call(
			merchant,
			"POST",
			"/v1/market/orders/"+o.ID+"/actions",
			map[string]any{"action": "ship", "carrier": "TEST", "tracking": "TEST-" + suffix},
			200,
		)
		call(buyer, "POST", "/v1/market/orders/"+o.ID+"/actions", map[string]any{"action": "receive"}, 200)
	}
	var dashboard struct {
		GMV     int64
		Payable int64
	}
	decode(call(buyer, "GET", "/v1/market/dashboard", nil, 200), &dashboard)
	if dashboard.GMV != 3980 || dashboard.Payable != 3782 {
		t.Fatalf("incorrect totals: %+v", dashboard)
	}
	call(adminToken, "PUT", "/v1/users/"+userIDs[2], map[string]any{"status": 0}, 200)
	call(buyer, "GET", "/v1/market/context", nil, 403)
	call("", "POST", "/login", map[string]any{"identifier": "registered_" + suffix, "password": password}, 401)
	t.Log(
		"real HTTP + JWT + gRPC + MySQL + Redis: two merchants, split checkout, replay, fulfillment, totals, tenant isolation and account disabling passed",
	)
}
