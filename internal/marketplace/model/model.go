// Package model 定义商城的数据模型及业务输入。
package model

import "time"

// Base 保存商城资源标识和时间戳。
type Base struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ID        string    `json:"id"        gorm:"primaryKey;size:36"`
}

// Shop 保存商家入驻资料和审核状态。
type Shop struct {
	Base
	OwnerID       string `json:"ownerId"       gorm:"size:64;uniqueIndex"`
	Name          string `json:"name"          gorm:"size:100"`
	Contact       string `json:"contact"       gorm:"size:100"`
	License       string `json:"license"       gorm:"size:200"`
	Status        string `json:"status"        gorm:"size:24;index"`
	Reason        string `json:"reason"        gorm:"size:500"`
	CommissionBPS int64  `json:"commissionBps"`
}

// Category 保存商品分类。
type Category struct {
	Base
	Name string `json:"name" gorm:"size:100;uniqueIndex"`
}

// Product 保存商品信息、审核状态及可用库存，价格单位为分。
type Product struct {
	Base
	ShopID      string `json:"shopId"      gorm:"size:36;index"`
	CategoryID  string `json:"categoryId"  gorm:"size:36;index"`
	Name        string `json:"name"        gorm:"size:200"`
	Description string `json:"description" gorm:"type:text"`
	Image       string `json:"image"       gorm:"size:1000"`
	Status      string `json:"status"      gorm:"size:24;index"`
	Reason      string `json:"reason"      gorm:"size:500"`
	Price       int64  `json:"price"`
	Stock       int64  `json:"stock"`
}

// Address 保存用户收货地址。
type Address struct {
	Base
	UserID    string `json:"userId"    gorm:"size:64;index"`
	Recipient string `json:"recipient" gorm:"size:100"`
	Phone     string `json:"phone"     gorm:"size:30"`
	Detail    string `json:"detail"    gorm:"size:500"`
}

// CartItem 保存用户选购的商品及数量。
type CartItem struct {
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Base
	UserID    string `json:"userId"            gorm:"size:64;uniqueIndex:cart_owner_product"`
	ProductID string `json:"productId"         gorm:"size:36;uniqueIndex:cart_owner_product"`
	Quantity  int64  `json:"quantity"`
}

// Checkout 记录结算幂等键和请求指纹。
type Checkout struct {
	Base
	UserID      string `json:"userId"     gorm:"size:64;uniqueIndex:checkout_request"`
	RequestKey  string `json:"requestKey" gorm:"size:100;uniqueIndex:checkout_request"`
	Fingerprint string `json:"-"          gorm:"size:64"`
}

// Order 保存按商家拆分的订单及交易快照，金额单位为分。
type Order struct {
	ExpiresAt time.Time `json:"expiresAt"`
	Base
	CheckoutID   string      `json:"checkoutId"   gorm:"size:36;index"`
	UserID       string      `json:"userId"       gorm:"size:64;index"`
	ShopID       string      `json:"shopId"       gorm:"size:36;index"`
	Status       string      `json:"status"       gorm:"size:24;index"`
	Recipient    string      `json:"recipient"    gorm:"size:100"`
	Phone        string      `json:"phone"        gorm:"size:30"`
	Address      string      `json:"address"      gorm:"size:500"`
	Carrier      string      `json:"carrier"      gorm:"size:100"`
	Tracking     string      `json:"tracking"     gorm:"size:100"`
	RefundReason string      `json:"refundReason" gorm:"size:500"`
	RefundFrom   string      `json:"refundFrom"   gorm:"size:24"`
	Items        []OrderItem `json:"items"        gorm:"foreignKey:OrderID"`
	Amount       int64       `json:"amount"`
	Commission   int64       `json:"commission"`
}

// OrderItem 保存下单时的商品名称、价格及数量。
type OrderItem struct {
	Base
	OrderID   string `json:"orderId"   gorm:"size:36;index"`
	ProductID string `json:"productId" gorm:"size:36;index"`
	Name      string `json:"name"      gorm:"size:200"`
	Price     int64  `json:"price"`
	Quantity  int64  `json:"quantity"`
}

// Audit 保存业务操作的审计记录。
type Audit struct {
	Base
	ActorID  string `json:"actorId"  gorm:"size:64;index"`
	Action   string `json:"action"   gorm:"size:100"`
	TargetID string `json:"targetId" gorm:"size:36;index"`
	Detail   string `json:"detail"   gorm:"size:1000"`
}

// Models 返回商城迁移涉及的全部数据模型。
func Models() []any {
	return []any{
		&Shop{},
		&Category{},
		&Product{},
		&Address{},
		&CartItem{},
		&Checkout{},
		&Order{},
		&OrderItem{},
		&Audit{},
	}
}
