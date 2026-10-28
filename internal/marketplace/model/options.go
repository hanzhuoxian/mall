package model

// Actor 标识当前操作人及其平台管理权限。
type Actor struct {
	ID    string
	Admin bool
}

// ListOptions 表示已归一化的分页参数。
type ListOptions struct {
	Page     int
	PageSize int
}

// Page 是商城列表接口的分页结果。
type Page[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}

// ProductFilter 限定商品查询范围。
type ProductFilter struct {
	OwnerID    string
	Name       string
	Status     string
	CategoryID string
	Catalog    bool
}

// OrderFilter 限定订单及汇总查询范围。
type OrderFilter struct {
	UserID  string
	OwnerID string
	Status  string
	DenyAll bool
}

// Summary 保存经营概览的金额汇总，金额单位为分。
type Summary struct {
	Orders     int64 `json:"orders"`
	GMV        int64 `json:"gmv"`
	Commission int64 `json:"commission"`
	Payable    int64 `json:"payable"`
}

// Context 返回用户可访问的商家及支付配置。
type Context struct {
	UserID         string `json:"userId"`
	Shops          []Shop `json:"shops"`
	Admin          bool   `json:"admin"`
	SandboxPayment bool   `json:"sandboxPayment"`
}

// ShopInput 是商家入驻申请。
type ShopInput struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
	License string `json:"license"`
}

// ShopReview 是平台对商家资质与佣金的审核结果。
type ShopReview struct {
	Status        string `json:"status"`
	Reason        string `json:"reason"`
	CommissionBPS int64  `json:"commissionBps"`
}

// ProductReview 是商品审核或下架请求。
type ProductReview struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// CheckoutInput 是购物车结算请求，RequestKey 用于幂等重试。
type CheckoutInput struct {
	AddressID  string   `json:"addressId"`
	RequestKey string   `json:"requestKey"`
	CartIDs    []string `json:"cartIds"`
}

// OrderAction 是订单操作及对应的物流、售后信息。
type OrderAction struct {
	Action   string `json:"action"`
	Carrier  string `json:"carrier"`
	Tracking string `json:"tracking"`
	Reason   string `json:"reason"`
}
