package contracts

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetMemberOrderPagingRequest(t *testing.T) {
	request := GetMemberOrderPagingRequest{
		MemberID:  "member-123",
		PageIndex: 1,
		PageSize:  10,
	}

	assert.Equal(t, "member-123", request.MemberID)
	assert.Equal(t, 1, request.PageIndex)
	assert.Equal(t, 10, request.PageSize)
}

func TestGetMemberOrderPagingResponse(t *testing.T) {
	now := time.Now()
	response := GetMemberOrderPagingResponse{
		ID:            "order-123",
		OrderNo:       "ORD202508120001",
		ProductName:   "拿铁咖啡",
		PayAmount:     decimal.NewFromFloat(15.80),
		CreatedAt:     now,
		PaymentStatus: PaymentStatusPaid,
	}

	assert.Equal(t, "order-123", response.ID)
	assert.Equal(t, "ORD202508120001", response.OrderNo)
	assert.Equal(t, "拿铁咖啡", response.ProductName)
	assert.True(t, response.PayAmount.Equal(decimal.NewFromFloat(15.80)))
	assert.Equal(t, now, response.CreatedAt)
	assert.Equal(t, PaymentStatusPaid, response.PaymentStatus)
}

func TestCreateOrderRequest(t *testing.T) {
	request := CreateOrderRequest{
		MemberID:  "member-123",
		MachineID: "machine-001",
		ProductID: "product-001",
		HasCup:    true,
		PayAmount: decimal.NewFromFloat(15.80),
	}

	assert.Equal(t, "member-123", request.MemberID)
	assert.Equal(t, "machine-001", request.MachineID)
	assert.Equal(t, "product-001", request.ProductID)
	assert.True(t, request.HasCup)
	assert.True(t, request.PayAmount.Equal(decimal.NewFromFloat(15.80)))
}

func TestCreateOrderResponse(t *testing.T) {
	response := CreateOrderResponse{
		OrderID: "order-123",
		OrderNo: "ORD202508120001",
		Message: "订单创建成功",
	}

	assert.Equal(t, "order-123", response.OrderID)
	assert.Equal(t, "ORD202508120001", response.OrderNo)
	assert.Equal(t, "订单创建成功", response.Message)
}

func TestRefundOrderRequest(t *testing.T) {
	request := RefundOrderRequest{
		OrderID:        "order-123",
		Reason:         "设备故障",
		IsMachineOwner: true,
	}

	assert.Equal(t, "order-123", request.OrderID)
	assert.Equal(t, "设备故障", request.Reason)
	assert.True(t, request.IsMachineOwner)
}

func TestRefundOrderResponse(t *testing.T) {
	response := RefundOrderResponse{
		OrderID:      "order-123",
		RefundAmount: decimal.NewFromFloat(15.80),
		Message:      "退款成功",
	}

	assert.Equal(t, "order-123", response.OrderID)
	assert.True(t, response.RefundAmount.Equal(decimal.NewFromFloat(15.80)))
	assert.Equal(t, "退款成功", response.Message)
}

func TestGetOrderByIdResponse(t *testing.T) {
	now := time.Now()
	paymentTime := now.Add(time.Minute * 5)
	refundReason := "设备故障"

	response := GetOrderByIdResponse{
		ID:            "order-123",
		OrderNo:       "ORD202508120001",
		MachineID:     "machine-001",
		MachineName:   "办公楼1层咖啡机",
		ProductID:     "product-001",
		ProductName:   "拿铁咖啡",
		PayAmount:     decimal.NewFromFloat(15.80),
		PaymentStatus: PaymentStatusPaid,
		MakeStatus:    MakeStatusMade,
		CreatedAt:     now,
		PaymentTime:   &paymentTime,
		HasCup:        true,
		RefundAmount:  decimal.NewFromFloat(0),
		RefundReason:  &refundReason,
	}

	assert.Equal(t, "order-123", response.ID)
	assert.Equal(t, "ORD202508120001", response.OrderNo)
	assert.Equal(t, "machine-001", response.MachineID)
	assert.Equal(t, "办公楼1层咖啡机", response.MachineName)
	assert.Equal(t, "product-001", response.ProductID)
	assert.Equal(t, "拿铁咖啡", response.ProductName)
	assert.True(t, response.PayAmount.Equal(decimal.NewFromFloat(15.80)))
	assert.Equal(t, PaymentStatusPaid, response.PaymentStatus)
	assert.Equal(t, MakeStatusMade, response.MakeStatus)
	assert.Equal(t, now, response.CreatedAt)
	assert.NotNil(t, response.PaymentTime)
	assert.Equal(t, paymentTime, *response.PaymentTime)
	assert.True(t, response.HasCup)
	assert.True(t, response.RefundAmount.Equal(decimal.NewFromFloat(0)))
	assert.NotNil(t, response.RefundReason)
	assert.Equal(t, "设备故障", *response.RefundReason)
}

func TestOrderPagingResponse(t *testing.T) {
	now := time.Now()
	orders := []GetMemberOrderPagingResponse{
		{
			ID:            "order-123",
			OrderNo:       "ORD202508120001",
			ProductName:   "拿铁咖啡",
			PayAmount:     decimal.NewFromFloat(15.80),
			CreatedAt:     now,
			PaymentStatus: PaymentStatusPaid,
		},
	}

	meta := PaginationMeta{
		Total:       1,
		Count:       1,
		PerPage:     10,
		CurrentPage: 1,
		TotalPages:  1,
		HasNext:     false,
		HasPrev:     false,
		Meta: &Meta{
			Timestamp: now,
			Version:   "v1.0.0",
		},
	}

	response := OrderPagingResponse{
		Orders: orders,
		Meta:   meta,
	}

	assert.Len(t, response.Orders, 1)
	assert.Equal(t, "order-123", response.Orders[0].ID)
	assert.Equal(t, int64(1), response.Meta.Total)
	assert.Equal(t, 1, response.Meta.Count)
}

func TestOrderStatusConstants(t *testing.T) {
	// 测试支付状态常量
	assert.Equal(t, "WaitPay", PaymentStatusWaitPay)
	assert.Equal(t, "Paid", PaymentStatusPaid)
	assert.Equal(t, "Refunded", PaymentStatusRefunded)
	assert.Equal(t, "Cancelled", PaymentStatusCancelled)

	// 测试制作状态常量
	assert.Equal(t, "WaitMake", MakeStatusWaitMake)
	assert.Equal(t, "Making", MakeStatusMaking)
	assert.Equal(t, "Made", MakeStatusMade)
	assert.Equal(t, "Failed", MakeStatusFailed)
}

func TestOrderErrorCodes(t *testing.T) {
	assert.Equal(t, "ORDER_NOT_FOUND", ErrorCodeOrderNotFound)
	assert.Equal(t, "ORDER_PERMISSION_DENIED", ErrorCodeOrderPermissionDenied)
	assert.Equal(t, "ORDER_ALREADY_PAID", ErrorCodeOrderAlreadyPaid)
	assert.Equal(t, "ORDER_ALREADY_REFUNDED", ErrorCodeOrderAlreadyRefunded)
	assert.Equal(t, "DEVICE_OFFLINE", ErrorCodeDeviceOffline)
	assert.Equal(t, "INSUFFICIENT_INVENTORY", ErrorCodeInsufficientInventory)
	assert.Equal(t, "ORDER_CREATE_FAILED", ErrorCodeOrderCreateFailed)
	assert.Equal(t, "REFUND_FAILED", ErrorCodeRefundFailed)
	assert.Equal(t, "INVALID_ORDER_STATUS", ErrorCodeInvalidOrderStatus)
	assert.Equal(t, "MACHINE_NOT_AVAILABLE", ErrorCodeMachineNotAvailable)
	assert.Equal(t, "PRODUCT_NOT_AVAILABLE", ErrorCodeProductNotAvailable)
}

func TestDecimalSerialization(t *testing.T) {
	request := CreateOrderRequest{
		MemberID:  "member-123",
		MachineID: "machine-001",
		ProductID: "product-001",
		HasCup:    true,
		PayAmount: decimal.NewFromFloat(15.80),
	}

	// 测试JSON序列化
	data, err := json.Marshal(request)
	assert.NoError(t, err)

	// 测试JSON反序列化
	var unmarshaled CreateOrderRequest
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, request.MemberID, unmarshaled.MemberID)
	assert.Equal(t, request.MachineID, unmarshaled.MachineID)
	assert.Equal(t, request.ProductID, unmarshaled.ProductID)
	assert.Equal(t, request.HasCup, unmarshaled.HasCup)
	assert.True(t, request.PayAmount.Equal(unmarshaled.PayAmount))
}
