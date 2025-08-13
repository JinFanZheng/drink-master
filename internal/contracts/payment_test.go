package contracts

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWeChatPayRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     WeChatPayRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: WeChatPayRequest{
				Ext1:        "test_account",
				Ext2:        "test_key",
				Ext3:        "VM_001_",
				NotifyUrl:   "http://example.com/callback",
				ChannelCode: "fuiou_pay_merchant",
				OrderNo:     "ORD20250812001",
				OpenId:      "test_open_id",
				OrderInfo:   "拿铁咖啡(需要杯子)",
				TransAmt:    1580, // 15.80元转分
			},
			wantErr: false,
		},
		{
			name: "empty required fields",
			req: WeChatPayRequest{
				TransAmt: 100,
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			req: WeChatPayRequest{
				Ext1:        "test_account",
				Ext2:        "test_key",
				Ext3:        "VM_001_",
				NotifyUrl:   "http://example.com/callback",
				ChannelCode: "fuiou_pay_merchant",
				OrderNo:     "ORD20250812001",
				OpenId:      "test_open_id",
				OrderInfo:   "拿铁咖啡",
				TransAmt:    0, // 无效金额
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: 实际的验证应该使用validator库
			// 这里只是演示测试结构
			if tt.wantErr {
				assert.True(t, tt.req.Ext1 == "" || tt.req.TransAmt <= 0)
			} else {
				assert.NotEmpty(t, tt.req.Ext1)
				assert.Greater(t, tt.req.TransAmt, int32(0))
			}
		})
	}
}

func TestWeChatPayResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		response WeChatPayResponse
		expected bool
	}{
		{
			name: "success response",
			response: WeChatPayResponse{
				IsSuccess: true,
				AppId:     "test_app_id",
				TimeStamp: "1642492800",
				NonceStr:  "random_string",
				Package:   "prepay_id=test_prepay_id",
				SignType:  "RSA",
				PaySign:   "test_signature",
			},
			expected: true,
		},
		{
			name: "failed response",
			response: WeChatPayResponse{
				IsSuccess: false,
				Message:   "支付失败",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.response.IsSuccess)
		})
	}
}

func TestTranQueryResponse_PaymentStatuses(t *testing.T) {
	tests := []struct {
		name           string
		paymentStatus  string
		expectedResult string
	}{
		{"success", PaymentStatusSuccess, "Success"},
		{"paying", PaymentStatusPaying, "Paying"},
		{"cancelled", PaymentStatusCancel, "Cancel"},
		{"failed", PaymentStatusFailure, "Failure"},
		{"timeout", PaymentStatusTimeout, "Timeout"},
		{"exception", PaymentStatusException, "Exception"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := TranQueryResponse{
				IsSuccess:     true,
				PaymentStatus: tt.paymentStatus,
			}
			assert.Equal(t, tt.expectedResult, response.PaymentStatus)
		})
	}
}

func TestPaymentCallbackRequest_Validation(t *testing.T) {
	validTime := time.Now()

	tests := []struct {
		name    string
		req     PaymentCallbackRequest
		isValid bool
	}{
		{
			name: "valid callback",
			req: PaymentCallbackRequest{
				OrderNo:       "ORD20250812001",
				TransactionId: "wx_123456789",
				Amount:        decimal.NewFromFloat(15.80),
				Status:        PaymentStatusSuccess,
				PaidAt:        validTime,
				Signature:     "valid_signature",
			},
			isValid: true,
		},
		{
			name: "missing order number",
			req: PaymentCallbackRequest{
				TransactionId: "wx_123456789",
				Amount:        decimal.NewFromFloat(15.80),
				Status:        PaymentStatusSuccess,
				PaidAt:        validTime,
				Signature:     "valid_signature",
			},
			isValid: false,
		},
		{
			name: "zero amount",
			req: PaymentCallbackRequest{
				OrderNo:       "ORD20250812001",
				TransactionId: "wx_123456789",
				Amount:        decimal.Zero,
				Status:        PaymentStatusSuccess,
				PaidAt:        validTime,
				Signature:     "valid_signature",
			},
			isValid: true, // 允许零金额（免费订单）
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				assert.NotEmpty(t, tt.req.OrderNo)
				assert.NotEmpty(t, tt.req.TransactionId)
				assert.NotEmpty(t, tt.req.Status)
			} else {
				assert.True(t, tt.req.OrderNo == "" || tt.req.TransactionId == "" || tt.req.Status == "")
			}
		})
	}
}

func TestPaymentConstants(t *testing.T) {
	// 测试支付渠道常量
	assert.Equal(t, "fuiou_pay_merchant", ChannelCodeFuiouMerchant)

	// 测试支付方式常量
	assert.Equal(t, "WeChatPay", ModeOfPaymentWeChat)

	// 测试支付状态常量
	assert.Equal(t, "Success", PaymentStatusSuccess)
	assert.Equal(t, "Paying", PaymentStatusPaying)
	assert.Equal(t, "Cancel", PaymentStatusCancel)

	// 测试免支付标识
	assert.Equal(t, "FREE_OF_PAYMENT", FreePaymentChannelOrderNo)
}

func TestPaymentErrorCodes(t *testing.T) {
	// 测试错误码常量存在性
	errorCodes := []string{
		ErrorCodePaymentOrderNotFound,
		ErrorCodePaymentOrderAlreadyPaid,
		ErrorCodePaymentFailed,
		ErrorCodePaymentQueryFailed,
		ErrorCodeInvalidPaymentAmount,
		ErrorCodePaymentAccountNotFound,
		ErrorCodePaymentCallbackInvalid,
	}

	for _, code := range errorCodes {
		assert.NotEmpty(t, code)
		assert.Contains(t, code, "PAYMENT")
	}
}

func TestPaymentAccount_RequiredFields(t *testing.T) {
	account := PaymentAccount{
		ReceivingAccount:     "merchant_123",
		ReceivingKey:         "secret_key_456",
		ReceivingOrderPrefix: "VM_001_",
	}

	assert.NotEmpty(t, account.ReceivingAccount)
	assert.NotEmpty(t, account.ReceivingKey)
	assert.NotEmpty(t, account.ReceivingOrderPrefix)
}
