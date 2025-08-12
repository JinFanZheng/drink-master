package enums

import "testing"

func TestPaymentStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   PaymentStatus
		expected int
	}{
		{"PaymentStatusWaitPay", PaymentStatusWaitPay, 0},
		{"PaymentStatusPaid", PaymentStatusPaid, 1},
		{"PaymentStatusInvalid", PaymentStatusInvalid, 2},
		{"PaymentStatusRefunded", PaymentStatusRefunded, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.status) != tt.expected {
				t.Errorf("Expected %s to be %d, but got %d", tt.name, tt.expected, int(tt.status))
			}
		})
	}
}

func TestGetPaymentStatusDesc(t *testing.T) {
	tests := []struct {
		name     string
		status   PaymentStatus
		expected string
	}{
		{"Wait pay status", PaymentStatusWaitPay, "待支付"},
		{"Paid status", PaymentStatusPaid, "已支付"},
		{"Invalid status", PaymentStatusInvalid, "已失效"},
		{"Refunded status", PaymentStatusRefunded, "已退款"},
		{"Unknown status", PaymentStatus(999), "未知状态"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPaymentStatusDesc(tt.status)
			if result != tt.expected {
				t.Errorf("Expected GetPaymentStatusDesc(%d) to be '%s', but got '%s'", tt.status, tt.expected, result)
			}
		})
	}
}

func TestPaymentStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   PaymentStatus
		expected string
	}{
		{"Wait pay status string", PaymentStatusWaitPay, "待支付"},
		{"Paid status string", PaymentStatusPaid, "已支付"},
		{"Invalid status string", PaymentStatusInvalid, "已失效"},
		{"Refunded status string", PaymentStatusRefunded, "已退款"},
		{"Unknown status string", PaymentStatus(999), "未知状态"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("Expected %d.String() to be '%s', but got '%s'", tt.status, tt.expected, result)
			}
		})
	}
}

func TestPaymentStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   PaymentStatus
		expected bool
	}{
		{"Valid WaitPay status", PaymentStatusWaitPay, true},
		{"Valid Paid status", PaymentStatusPaid, true},
		{"Valid Invalid status", PaymentStatusInvalid, true},
		{"Valid Refunded status", PaymentStatusRefunded, true},
		{"Invalid status -1", PaymentStatus(-1), false},
		{"Invalid status 4", PaymentStatus(4), false},
		{"Invalid large status", PaymentStatus(999), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			if result != tt.expected {
				t.Errorf("Expected %d.IsValid() to be %t, but got %t", tt.status, tt.expected, result)
			}
		})
	}
}
