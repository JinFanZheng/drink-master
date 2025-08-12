package contracts

import (
	"time"

	"github.com/shopspring/decimal"
)

// GetMemberOrderPagingRequest 获取会员订单分页列表请求
type GetMemberOrderPagingRequest struct {
	MemberID  string `json:"memberId"`
	PageIndex int    `json:"pageIndex" validate:"min=1" example:"1"`
	PageSize  int    `json:"pageSize" validate:"min=1,max=100" example:"10"`
}

// GetMemberOrderPagingResponse 会员订单分页列表响应
type GetMemberOrderPagingResponse struct {
	ID            string          `json:"id" example:"order-123"`
	OrderNo       string          `json:"orderNo" example:"ORD202508120001"`
	ProductName   string          `json:"productName" example:"拿铁咖啡"`
	PayAmount     decimal.Decimal `json:"payAmount" example:"15.80"`
	CreatedAt     time.Time       `json:"createdAt" example:"2025-08-12T10:30:00Z"`
	PaymentStatus string          `json:"paymentStatus" example:"Paid"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	MemberID  string          `json:"memberId"`
	MachineID string          `json:"machineId" binding:"required" validate:"required" example:"machine-001"`
	ProductID string          `json:"productId" binding:"required" validate:"required" example:"product-001"`
	HasCup    bool            `json:"hasCup" example:"true"`
	PayAmount decimal.Decimal `json:"payAmount" binding:"required" validate:"required,gt=0" example:"15.80"`
}

// CreateOrderResponse 创建订单响应
type CreateOrderResponse struct {
	OrderID string `json:"orderId" example:"order-123"`
	OrderNo string `json:"orderNo" example:"ORD202508120001"`
	Message string `json:"message" example:"订单创建成功"`
}

// RefundOrderRequest 退款订单请求
type RefundOrderRequest struct {
	OrderID        string `json:"orderId" binding:"required" validate:"required" example:"order-123"`
	Reason         string `json:"reason" example:"设备故障无法出货"`
	IsMachineOwner bool   `json:"isMachineOwner"`
}

// RefundOrderResponse 退款订单响应
type RefundOrderResponse struct {
	OrderID      string          `json:"orderId" example:"order-123"`
	RefundAmount decimal.Decimal `json:"refundAmount" example:"15.80"`
	Message      string          `json:"message" example:"退款成功"`
}

// GetOrderByIdResponse 根据ID获取订单详情响应
type GetOrderByIdResponse struct {
	ID            string          `json:"id" example:"order-123"`
	OrderNo       string          `json:"orderNo" example:"ORD202508120001"`
	MachineID     string          `json:"machineId" example:"machine-001"`
	MachineName   string          `json:"machineName" example:"办公楼1层咖啡机"`
	ProductID     string          `json:"productId" example:"product-001"`
	ProductName   string          `json:"productName" example:"拿铁咖啡"`
	PayAmount     decimal.Decimal `json:"payAmount" example:"15.80"`
	PaymentStatus string          `json:"paymentStatus" example:"Paid"`
	MakeStatus    string          `json:"makeStatus" example:"Made"`
	CreatedAt     time.Time       `json:"createdAt" example:"2025-08-12T10:30:00Z"`
	PaymentTime   *time.Time      `json:"paymentTime,omitempty" example:"2025-08-12T10:30:30Z"`
	HasCup        bool            `json:"hasCup" example:"true"`
	RefundAmount  decimal.Decimal `json:"refundAmount" example:"0"`
	RefundReason  *string         `json:"refundReason,omitempty"`
}

// OrderPagingResponse 订单分页响应
type OrderPagingResponse struct {
	Orders []GetMemberOrderPagingResponse `json:"orders"`
	Meta   PaginationMeta                 `json:"meta"`
}

// 订单状态常量
const (
	PaymentStatusWaitPay   = "WaitPay"   // 等待支付
	PaymentStatusPaid      = "Paid"      // 已支付
	PaymentStatusRefunded  = "Refunded"  // 已退款
	PaymentStatusCancelled = "Cancelled" // 已取消

	MakeStatusWaitMake = "WaitMake" // 等待制作
	MakeStatusMaking   = "Making"   // 制作中
	MakeStatusMade     = "Made"     // 已完成
	MakeStatusFailed   = "Failed"   // 制作失败
)

// 订单错误码常量
const (
	ErrorCodeOrderNotFound         = "ORDER_NOT_FOUND"
	ErrorCodeOrderPermissionDenied = "ORDER_PERMISSION_DENIED"
	ErrorCodeOrderAlreadyPaid      = "ORDER_ALREADY_PAID"
	ErrorCodeOrderAlreadyRefunded  = "ORDER_ALREADY_REFUNDED"
	ErrorCodeDeviceOffline         = "DEVICE_OFFLINE"
	ErrorCodeInsufficientInventory = "INSUFFICIENT_INVENTORY"
	ErrorCodeOrderCreateFailed     = "ORDER_CREATE_FAILED"
	ErrorCodeRefundFailed          = "REFUND_FAILED"
	ErrorCodeInvalidOrderStatus    = "INVALID_ORDER_STATUS"
	ErrorCodeMachineNotAvailable   = "MACHINE_NOT_AVAILABLE"
	ErrorCodeProductNotAvailable   = "PRODUCT_NOT_AVAILABLE"
)
