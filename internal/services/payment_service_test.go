package services

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

// MockOrderRepository for testing
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(id string) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByMemberPaging(memberID string, pageIndex, pageSize int) ([]models.Order, int64, error) {
	args := m.Called(memberID, pageIndex, pageSize)
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) Update(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByOrderNo(orderNo string) (*models.Order, error) {
	args := m.Called(orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

// 使用已存在的MockMachineRepository

func TestPaymentService_WeChatPay_Mock(t *testing.T) {
	// 设置Mock模式
	os.Setenv("MOCK_MODE", "true")
	defer os.Unsetenv("MOCK_MODE")

	mockOrderRepo := &MockOrderRepository{}
	mockMachineRepo := &MockMachineRepository{}

	service := &paymentService{
		orderRepo:   mockOrderRepo,
		machineRepo: mockMachineRepo,
	}

	req := contracts.WeChatPayRequest{
		Ext1:        "test_account",
		Ext2:        "test_key",
		Ext3:        "VM_001_",
		NotifyUrl:   "http://example.com/callback",
		ChannelCode: contracts.ChannelCodeFuiouMerchant,
		OrderNo:     "ORD20250812001",
		OpenId:      "test_open_id",
		OrderInfo:   "拿铁咖啡",
		TransAmt:    1580,
	}

	response, err := service.WeChatPay(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.IsSuccess)
	assert.Equal(t, "mock_app_id", response.AppId)
	assert.NotEmpty(t, response.TimeStamp)
	assert.Equal(t, "mock_nonce_str", response.NonceStr)
}

func TestPaymentService_TranQuery_Mock(t *testing.T) {
	// 设置Mock模式
	os.Setenv("MOCK_MODE", "true")
	defer os.Unsetenv("MOCK_MODE")

	mockOrderRepo := &MockOrderRepository{}
	mockMachineRepo := &MockMachineRepository{}

	service := &paymentService{
		orderRepo:   mockOrderRepo,
		machineRepo: mockMachineRepo,
	}

	req := contracts.TranQueryRequest{
		Ext1:          "test_account",
		Ext2:          "test_key",
		Ext3:          "VM_001_",
		ChannelCode:   contracts.ChannelCodeFuiouMerchant,
		OrderNo:       "ORD20250812001",
		ModeOfPayment: contracts.ModeOfPaymentWeChat,
	}

	response, err := service.TranQuery(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.IsSuccess)
	assert.Equal(t, contracts.PaymentStatusSuccess, response.PaymentStatus)
	assert.Equal(t, "mock_transaction_id_ORD20250812001", response.TransactionId)
}

func TestPaymentService_GetPaymentAccount(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockMachineRepo := &MockMachineRepository{}

	service := &paymentService{
		orderRepo:   mockOrderRepo,
		machineRepo: mockMachineRepo,
	}

	machine := &models.Machine{
		ID:        "machine-001",
		MachineNo: stringPtr("VM001"),
		Name:      stringPtr("测试咖啡机"),
	}

	mockMachineRepo.On("GetByID", "machine-001").Return(machine, nil)

	account, err := service.GetPaymentAccount("machine-001")

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.NotEmpty(t, account.ReceivingAccount)
	assert.NotEmpty(t, account.ReceivingKey)
	assert.Equal(t, "VM_VM001_", account.ReceivingOrderPrefix)

	mockMachineRepo.AssertExpectations(t)
}

func TestPaymentService_PayOrder(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockMachineRepo := &MockMachineRepository{}

	service := &paymentService{
		orderRepo:   mockOrderRepo,
		machineRepo: mockMachineRepo,
	}

	order := &models.Order{
		ID:            "order-001",
		PaymentStatus: int(enums.PaymentStatusWaitPay),
		PayAmount:     15.80,
	}

	mockOrderRepo.On("GetByID", "order-001").Return(order, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType("*models.Order")).Return(nil)

	paidTime := time.Now()
	req := contracts.PayOrderRequest{
		ID:             "order-001",
		ChannelOrderNo: "wx_123456789",
		PaidAt:         paidTime,
	}

	err := service.PayOrder(req)

	assert.NoError(t, err)

	// 验证订单状态已更新
	mockOrderRepo.AssertCalled(t, "Update", mock.MatchedBy(func(o *models.Order) bool {
		return o.PaymentStatus == int(enums.PaymentStatusPaid) &&
			o.PaymentTime != nil &&
			*o.ChannelOrderNo == "wx_123456789"
	}))

	mockOrderRepo.AssertExpectations(t)
}

func TestPaymentService_PayOrder_InvalidStatus(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockMachineRepo := &MockMachineRepository{}

	service := &paymentService{
		orderRepo:   mockOrderRepo,
		machineRepo: mockMachineRepo,
	}

	order := &models.Order{
		ID:            "order-001",
		PaymentStatus: int(enums.PaymentStatusPaid), // 已经是已支付状态
		PayAmount:     15.80,
	}

	mockOrderRepo.On("GetByID", "order-001").Return(order, nil)

	req := contracts.PayOrderRequest{
		ID:             "order-001",
		ChannelOrderNo: "wx_123456789",
		PaidAt:         time.Now(),
	}

	err := service.PayOrder(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in wait pay status")

	mockOrderRepo.AssertExpectations(t)
}

func TestPaymentService_InvalidOrder(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockMachineRepo := &MockMachineRepository{}

	service := &paymentService{
		orderRepo:   mockOrderRepo,
		machineRepo: mockMachineRepo,
	}

	order := &models.Order{
		ID:            "order-001",
		PaymentStatus: int(enums.PaymentStatusWaitPay),
		PayAmount:     15.80,
	}

	mockOrderRepo.On("GetByID", "order-001").Return(order, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType("*models.Order")).Return(nil)

	req := contracts.InvalidOrderRequest{
		ID: "order-001",
	}

	err := service.InvalidOrder(req)

	assert.NoError(t, err)

	// 验证订单状态已更新为失效
	mockOrderRepo.AssertCalled(t, "Update", mock.MatchedBy(func(o *models.Order) bool {
		return o.PaymentStatus == int(enums.PaymentStatusInvalid)
	}))

	mockOrderRepo.AssertExpectations(t)
}

func TestPaymentService_ProcessPaymentCallback_Success(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockMachineRepo := &MockMachineRepository{}

	service := &paymentService{
		orderRepo:   mockOrderRepo,
		machineRepo: mockMachineRepo,
	}

	order := &models.Order{
		ID:            "order-001",
		OrderNo:       stringPtr("ORD20250812001"),
		PaymentStatus: int(enums.PaymentStatusWaitPay),
		PayAmount:     15.80,
	}

	mockOrderRepo.On("GetByOrderNo", "ORD20250812001").Return(order, nil)
	mockOrderRepo.On("GetByID", "order-001").Return(order, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType("*models.Order")).Return(nil)

	req := contracts.PaymentCallbackRequest{
		OrderNo:       "ORD20250812001",
		TransactionId: "wx_123456789",
		Amount:        decimal.NewFromFloat(15.80),
		Status:        contracts.PaymentStatusSuccess,
		PaidAt:        time.Now(),
		Signature:     "valid_signature",
	}

	response, err := service.ProcessPaymentCallback(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Processed)
	assert.Equal(t, "回调处理成功", response.Message)

	mockOrderRepo.AssertExpectations(t)
}

func TestPaymentService_ProcessPaymentCallback_Failed(t *testing.T) {
	mockOrderRepo := &MockOrderRepository{}
	mockMachineRepo := &MockMachineRepository{}

	service := &paymentService{
		orderRepo:   mockOrderRepo,
		machineRepo: mockMachineRepo,
	}

	order := &models.Order{
		ID:            "order-001",
		OrderNo:       stringPtr("ORD20250812001"),
		PaymentStatus: int(enums.PaymentStatusWaitPay),
		PayAmount:     15.80,
	}

	mockOrderRepo.On("GetByOrderNo", "ORD20250812001").Return(order, nil)
	mockOrderRepo.On("GetByID", "order-001").Return(order, nil)
	mockOrderRepo.On("Update", mock.AnythingOfType("*models.Order")).Return(nil)

	req := contracts.PaymentCallbackRequest{
		OrderNo:       "ORD20250812001",
		TransactionId: "wx_123456789",
		Amount:        decimal.NewFromFloat(15.80),
		Status:        contracts.PaymentStatusFailure,
		PaidAt:        time.Now(),
		Signature:     "valid_signature",
	}

	response, err := service.ProcessPaymentCallback(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Processed)

	// 验证订单被作废
	mockOrderRepo.AssertCalled(t, "Update", mock.MatchedBy(func(o *models.Order) bool {
		return o.PaymentStatus == int(enums.PaymentStatusInvalid)
	}))

	mockOrderRepo.AssertExpectations(t)
}

func TestGetEnvOrDefault(t *testing.T) {
	// 测试存在的环境变量
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	result := getEnvOrDefault("TEST_ENV_VAR", "default_value")
	assert.Equal(t, "test_value", result)

	// 测试不存在的环境变量
	result = getEnvOrDefault("NON_EXISTENT_VAR", "default_value")
	assert.Equal(t, "default_value", result)
}
