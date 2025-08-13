package contracts

import (
	"time"

	"github.com/shopspring/decimal"
)

// WeChatPayRequest 微信支付请求
type WeChatPayRequest struct {
	Ext1        string `json:"ext1" validate:"required"`        // 收款账户
	Ext2        string `json:"ext2" validate:"required"`        // 收款密钥
	Ext3        string `json:"ext3" validate:"required"`        // 订单前缀
	NotifyUrl   string `json:"notifyUrl" validate:"required"`   // 回调地址
	ChannelCode string `json:"channelCode" validate:"required"` // 渠道代码
	OrderNo     string `json:"orderNo" validate:"required"`     // 订单号
	OpenId      string `json:"openId" validate:"required"`      // 微信OpenId
	Attach      string `json:"attach"`                          // 附加信息
	OrderInfo   string `json:"orderInfo" validate:"required"`   // 订单信息
	TransAmt    int32  `json:"transAmt" validate:"gt=0"`        // 交易金额(分)
}

// WeChatPayResponse 微信支付响应
type WeChatPayResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	AppId     string `json:"appId,omitempty"`
	TimeStamp string `json:"timeStamp,omitempty"`
	NonceStr  string `json:"nonceStr,omitempty"`
	Package   string `json:"package,omitempty"`
	SignType  string `json:"signType,omitempty"`
	PaySign   string `json:"paySign,omitempty"`
	Message   string `json:"message,omitempty"`
}

// TranQueryRequest 支付查询请求
type TranQueryRequest struct {
	Ext1          string `json:"ext1" validate:"required"`          // 收款账户
	Ext2          string `json:"ext2" validate:"required"`          // 收款密钥
	Ext3          string `json:"ext3" validate:"required"`          // 订单前缀
	ChannelCode   string `json:"channelCode" validate:"required"`   // 渠道代码
	OrderNo       string `json:"orderNo" validate:"required"`       // 订单号
	ModeOfPayment string `json:"modeOfPayment" validate:"required"` // 支付方式
}

// TranQueryResponse 支付查询响应
type TranQueryResponse struct {
	IsSuccess     bool      `json:"isSuccess"`
	PaymentStatus string    `json:"paymentStatus,omitempty"` // Success, Paying, Cancel, Failure, Timeout, Exception
	TransactionId string    `json:"transactionId,omitempty"` // 第三方交易号
	PaymentTime   time.Time `json:"paymentTime,omitempty"`   // 支付时间
	Message       string    `json:"message,omitempty"`       // 响应消息
	ErrorCode     string    `json:"errorCode,omitempty"`     // 错误码
}

// PayOrderRequest 支付订单请求（内部服务调用）
type PayOrderRequest struct {
	ID             string    `json:"id" validate:"required"`             // 订单ID
	ChannelOrderNo string    `json:"channelOrderNo" validate:"required"` // 第三方订单号
	PaidAt         time.Time `json:"paidAt" validate:"required"`         // 支付时间
}

// InvalidOrderRequest 作废订单请求（内部服务调用）
type InvalidOrderRequest struct {
	ID string `json:"id" validate:"required"` // 订单ID
}

// PaymentAccount 支付账户信息
type PaymentAccount struct {
	ReceivingAccount     string `json:"receivingAccount" validate:"required"`     // 收款账户
	ReceivingKey         string `json:"receivingKey" validate:"required"`         // 收款密钥
	ReceivingOrderPrefix string `json:"receivingOrderPrefix" validate:"required"` // 订单前缀
}

// GetPaymentRequest 获取支付信息请求
type GetPaymentRequest struct {
	OrderID string `form:"orderId" validate:"required" example:"order-123"`
}

// GetPaymentResponse 获取支付信息响应
type GetPaymentResponse struct {
	Code    int                `json:"code" example:"200"`
	Message string             `json:"message" example:"获取支付信息成功"`
	Data    *WeChatPayResponse `json:"data,omitempty"`
}

// QueryPaymentRequest 查询支付状态请求
type QueryPaymentRequest struct {
	OrderID string `form:"orderId" validate:"required" example:"order-123"`
}

// QueryPaymentResponse 查询支付状态响应
type QueryPaymentResponse struct {
	Message string `json:"message" example:"支付成功"`
}

// PaymentCallbackRequest 支付回调请求
type PaymentCallbackRequest struct {
	OrderNo       string          `json:"orderNo" validate:"required"`       // 订单号
	TransactionId string          `json:"transactionId" validate:"required"` // 第三方交易号
	Amount        decimal.Decimal `json:"amount" validate:"required"`        // 支付金额
	Status        string          `json:"status" validate:"required"`        // 支付状态
	PaidAt        time.Time       `json:"paidAt" validate:"required"`        // 支付时间
	Signature     string          `json:"signature" validate:"required"`     // 签名
}

// PaymentCallbackResponse 支付回调响应
type PaymentCallbackResponse struct {
	Processed bool   `json:"processed" example:"true"`
	Message   string `json:"message" example:"回调处理成功"`
}

// PaymentCallbackResultRequest 支付结果回调请求（基于VendingMachine.MobileAPI）
type PaymentCallbackResultRequest struct {
	ChannelCode    string    `json:"channelCode" validate:"required"`    // 渠道编码
	TransAmt       int       `json:"transAmt" validate:"required,gt=0"`  // 订单金额(分)
	ReturnAmt      int       `json:"returnAmt"`                          // 返还金额(分)
	OrderNo        string    `json:"orderNo" validate:"required"`        // 订单号
	OrderInfo      string    `json:"orderInfo"`                          // 订单信息
	ModeOfPayment  int       `json:"modeOfPayment" validate:"required"`  // 支付方式
	ChannelOrderNo string    `json:"channelOrderNo" validate:"required"` // 渠道订单号
	PaymentTime    time.Time `json:"paymentTime" validate:"required"`    // 支付时间
	CallbackType   string    `json:"callbackType" validate:"required"`   // 回调类型
}

// 支付相关常量
const (
	// 支付渠道
	ChannelCodeFuiouMerchant = "fuiou_pay_merchant"

	// 支付方式
	ModeOfPaymentWeChat = "WeChatPay"

	// 支付状态（第三方返回）
	PaymentStatusPaying    = "Paying"    // 支付中
	PaymentStatusSuccess   = "Success"   // 支付成功
	PaymentStatusCancel    = "Cancel"    // 已取消
	PaymentStatusFailure   = "Failure"   // 支付失败
	PaymentStatusTimeout   = "Timeout"   // 支付超时
	PaymentStatusException = "Exception" // 异常

	// 免支付标识
	FreePaymentChannelOrderNo = "FREE_OF_PAYMENT"
)

// 支付错误码
const (
	ErrorCodePaymentOrderNotFound    = "PAYMENT_ORDER_NOT_FOUND"
	ErrorCodePaymentOrderAlreadyPaid = "PAYMENT_ORDER_ALREADY_PAID"
	ErrorCodePaymentFailed           = "PAYMENT_FAILED"
	ErrorCodePaymentQueryFailed      = "PAYMENT_QUERY_FAILED"
	ErrorCodeInvalidPaymentAmount    = "INVALID_PAYMENT_AMOUNT"
	ErrorCodePaymentAccountNotFound  = "PAYMENT_ACCOUNT_NOT_FOUND"
	ErrorCodePaymentCallbackInvalid  = "PAYMENT_CALLBACK_INVALID"
)
