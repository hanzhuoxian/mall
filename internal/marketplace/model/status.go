package model

const (
	MaxReasonLength        = 500
	MaxEntityNameLength    = 100
	MaxProductNameLength   = 200
	MaxContactLength       = 100
	MaxLicenseLength       = 200
	MaxAddressDetailLength = 500
	MaxCarrierLength       = 100
	MaxTrackingLength      = 100
	MaxPageSize            = 100
	MaxPhoneLength         = 30
	MaxCartQuantity        = 999
	MaxCheckoutItems       = 100
	MaxRequestKeyLength    = 100
	MaxProductPriceCents   = 100_000_000
	MaxProductStock        = 100_000_000
	MaxStockAdjustment     = 1_000_000
	CommissionBasisPoints  = 10_000
	DefaultCommissionBPS   = 500

	SearchOperationIndex  = "index"
	SearchOperationDelete = "delete"

	ShopStatusPending   = "pending"
	ShopStatusApproved  = "approved"
	ShopStatusRejected  = "rejected"
	ShopStatusSuspended = "suspended"

	ProductStatusPending  = "pending"
	ProductStatusApproved = "approved"
	ProductStatusRejected = "rejected"
	ProductStatusOffline  = "offline"

	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusShipped   = "shipped"
	OrderStatusCompleted = "completed"
	OrderStatusRefunding = "refunding"
	OrderStatusRefunded  = "refunded"
	OrderStatusRejected  = "rejected"

	OrderActionPay           = "pay"
	OrderActionCancel        = "cancel"
	OrderActionShip          = "ship"
	OrderActionReceive       = "receive"
	OrderActionRefund        = "refund"
	OrderActionApproveRefund = "approve-refund"
	OrderActionRejectRefund  = "reject-refund"
)

// OrderCanceled 是取消订单状态，保留既有数据库及接口拼写。
const OrderCanceled = "cancelled" //nolint:misspell // Existing database and API value.
