package mysql_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
	"github.com/hanzhuoxian/mall/internal/marketplace/store"
	mysqlstore "github.com/hanzhuoxian/mall/internal/marketplace/store/mysql"
)

func testStore(t *testing.T) store.Store {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(model.Models()...); err != nil {
		t.Fatal(err)
	}
	s := mysqlstore.NewStore(db)
	t.Cleanup(func() {
		if err := s.Close(); err != nil {
			t.Error(err)
		}
	})
	return s
}

func TestTransactionRollbackAndContext(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	rollback := errors.New("rollback")
	err := s.Transaction(ctx, func(tx store.Store) error {
		if err := tx.SaveShop(ctx, &model.Shop{Base: model.Base{ID: "shop"}, OwnerID: "seller"}); err != nil {
			return err
		}
		if _, err := tx.Shop(ctx, "shop", false); err != nil {
			return err
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatalf("transaction error = %v", err)
	}
	if _, err := s.Shop(ctx, "shop", false); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("rolled back shop error = %v", err)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := s.ListShops(
		canceled,
		"",
		model.ListOptions{Page: 1, PageSize: 10},
	); !errors.Is(
		err,
		context.Canceled,
	) {
		t.Fatalf("canceled query error = %v", err)
	}
}

func TestCatalogPaginationAndSummaryScopes(t *testing.T) {
	s := testStore(t)
	ctx := context.Background()
	for _, shop := range []model.Shop{
		{Base: model.Base{ID: "approved"}, OwnerID: "seller", Status: "approved"},
		{Base: model.Base{ID: "pending"}, OwnerID: "other", Status: "pending"},
	} {
		if err := s.SaveShop(ctx, &shop); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now()
	for _, p := range []model.Product{
		{Base: model.Base{ID: "a", CreatedAt: now}, ShopID: "approved", Status: "approved"},
		{Base: model.Base{ID: "b", CreatedAt: now}, ShopID: "approved", Status: "approved"},
		{Base: model.Base{ID: "hidden"}, ShopID: "pending", Status: "approved"},
		{Base: model.Base{ID: "draft"}, ShopID: "approved", Status: "pending"},
	} {
		if err := s.SaveProduct(ctx, &p); err != nil {
			t.Fatal(err)
		}
	}
	page, err := s.ListProducts(ctx, model.ProductFilter{Catalog: true}, model.ListOptions{Page: 2, PageSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 2 || len(page.Items) != 1 || page.Items[0].ID != "b" {
		t.Fatalf("catalog page = %+v", page)
	}
	for _, order := range []model.Order{
		{Base: model.Base{ID: "paid"}, ShopID: "approved", Status: "paid", Amount: 100, Commission: 10},
		{Base: model.Base{ID: "done"}, ShopID: "approved", Status: "completed", Amount: 200, Commission: 20},
		{Base: model.Base{ID: "unpaid"}, ShopID: "approved", Status: "pending", Amount: 400, Commission: 40},
		{Base: model.Base{ID: "other"}, ShopID: "pending", Status: "completed", Amount: 800, Commission: 80},
	} {
		if err = s.CreateOrder(ctx, &order); err != nil {
			t.Fatal(err)
		}
	}
	summary, err := s.Summary(ctx, model.OrderFilter{OwnerID: "seller"})
	if err != nil {
		t.Fatal(err)
	}
	if summary != (model.Summary{Orders: 3, GMV: 300, Commission: 30, Payable: 180}) {
		t.Fatalf("summary = %+v", summary)
	}
	denied, err := s.ListOrders(ctx, model.OrderFilter{DenyAll: true}, model.ListOptions{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if denied.Total != 0 || denied.Items == nil || len(denied.Items) != 0 {
		t.Fatalf("denied page = %+v", denied)
	}
}
