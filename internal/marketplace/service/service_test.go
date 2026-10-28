package service

import (
	"context"
	"errors"

	"strings"
	"testing"
	"time"

	"github.com/hanzhuoxian/mall/internal/marketplace/model"
)

func TestOrderTransitions(t *testing.T) {
	cases := []struct {
		name, status, action, refundFrom, want string
		buyer, seller, admin, sandbox          bool
	}{
		{name: "buyer pays", status: "pending", action: "pay", want: "paid", buyer: true, sandbox: true},
		{name: "production payment blocked", status: "pending", action: "pay", buyer: true},
		{name: "seller cannot pay buyer order", status: "pending", action: "pay", seller: true, sandbox: true},
		{name: "admin cannot impersonate buyer", status: "pending", action: "pay", admin: true, sandbox: true},
		{name: "cancel pending", status: "pending", action: "cancel", want: model.OrderCanceled, buyer: true},
		{name: "cancel paid blocked", status: "paid", action: "cancel", buyer: true},
		{name: "ship paid", status: "paid", action: "ship", want: "shipped", seller: true},
		{name: "buyer cannot ship", status: "paid", action: "ship", buyer: true},
		{name: "unpaid cannot ship", status: "pending", action: "ship", seller: true},
		{name: "confirm receipt", status: "shipped", action: "receive", want: "completed", buyer: true},
		{name: "seller cannot confirm", status: "shipped", action: "receive", seller: true},
		{name: "request refund", status: "completed", action: "refund", want: "refunding", buyer: true},
		{name: "unpaid refund denied", status: "pending", action: "refund", buyer: true},
		{
			name:    "approve sandbox refund",
			status:  "refunding",
			action:  "approve-refund",
			want:    "refunded",
			seller:  true,
			sandbox: true,
		},
		{name: "production refund denied", status: "refunding", action: "approve-refund", seller: true},
		{
			name:    "buyer cannot approve own refund",
			status:  "refunding",
			action:  "approve-refund",
			buyer:   true,
			sandbox: true,
		},
		{
			name:    "platform arbitrates refund",
			status:  "refunding",
			action:  "approve-refund",
			want:    "refunded",
			admin:   true,
			sandbox: true,
		},
		{
			name:       "reject restores previous status",
			status:     "refunding",
			action:     "reject-refund",
			refundFrom: "shipped",
			want:       "shipped",
			seller:     true,
		},
		{
			name:       "invalid refund origin",
			status:     "refunding",
			action:     "reject-refund",
			refundFrom: model.OrderCanceled,
			seller:     true,
		},
		{name: "double refund denied", status: "refunded", action: "approve-refund", seller: true, sandbox: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := model.Order{Status: tc.status, RefundFrom: tc.refundFrom, ExpiresAt: time.Now().Add(time.Hour)}
			got, err := transition(o, tc.action, tc.buyer, tc.seller, tc.admin, tc.sandbox)
			if tc.want == "" {
				if err == nil {
					t.Fatal("expected rejection")
				}
			} else if err != nil || got != tc.want {
				t.Fatalf("got %s %v", got, err)
			}
		})
	}
	if _, err := transition(
		model.Order{Status: "pending", ExpiresAt: time.Now().Add(-time.Minute)},
		"pay",
		true,
		false,
		false,
		true,
	); err == nil {
		t.Fatal("expired order paid")
	}
}
func TestCheckoutFingerprint(t *testing.T) {
	a := model.CheckoutInput{AddressID: "a", CartIDs: []string{"1", "2"}, RequestKey: "k"}
	b := model.CheckoutInput{AddressID: "a", CartIDs: []string{"2", "1"}, RequestKey: "k"}
	if fingerprint(a) != fingerprint(b) {
		t.Fatal("ordering should not affect retry")
	}
	b.AddressID = "other"
	if fingerprint(a) == fingerprint(b) {
		t.Fatal("address must affect fingerprint")
	}
	for _, in := range []model.CheckoutInput{{}, {RequestKey: "k"}, {RequestKey: "k", CartIDs: []string{"1", "1"}}} {
		if _, err := NewService(
			nil,
			false,
		).Checkout(context.Background(), model.Actor{ID: "u"}, in); !errors.Is(
			err,
			ErrInvalid,
		) {
			t.Fatal("invalid request accepted")
		}
	}
}
func TestProductValidation(t *testing.T) {
	p := model.Product{Name: "商品", Price: 1, Stock: 0}
	if !validProduct(p) {
		t.Fatal("valid product rejected")
	}
	for _, edit := range []func(*model.Product){func(p *model.Product) { p.Price = -1 }, func(p *model.Product) { p.Stock = -1 }, func(p *model.Product) { p.Price = 100000001 }, func(p *model.Product) { p.Name = "  " }, func(p *model.Product) { p.Image = "javascript:alert(1)" }, func(p *model.Product) { p.Description = strings.Repeat("字", 10001) }} {
		q := p
		edit(&q)
		if validProduct(q) {
			t.Fatal("invalid product accepted")
		}
	}
}
