package enums

// PaymentStatus represents the payment status of an order
type PaymentStatus int

const (
	// PaymentStatusWaitPay represents an order waiting for payment
	PaymentStatusWaitPay PaymentStatus = 0 // 待支付
	// PaymentStatusPaid represents an order that has been paid
	PaymentStatusPaid PaymentStatus = 1 // 已支付
	// PaymentStatusInvalid represents an order with invalid payment
	PaymentStatusInvalid PaymentStatus = 2 // 已失效
	// PaymentStatusRefunded represents an order that has been refunded
	PaymentStatusRefunded PaymentStatus = 3 // 已退款
)

// GetPaymentStatusDesc returns the description of the payment status
func GetPaymentStatusDesc(status PaymentStatus) string {
	switch status {
	case PaymentStatusWaitPay:
		return "待支付"
	case PaymentStatusPaid:
		return "已支付"
	case PaymentStatusInvalid:
		return "已失效"
	case PaymentStatusRefunded:
		return "已退款"
	default:
		return "未知状态"
	}
}

// String returns the string representation of the payment status
func (ps PaymentStatus) String() string {
	return GetPaymentStatusDesc(ps)
}

// IsValid checks if the payment status is valid
func (ps PaymentStatus) IsValid() bool {
	return ps >= PaymentStatusWaitPay && ps <= PaymentStatusRefunded
}
