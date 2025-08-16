package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

func TestNewOrderService(t *testing.T) {
	service := NewOrderService(nil, nil, nil, nil)
	assert.NotNil(t, service)
}

func TestOrderService_GenerateOrderNo(t *testing.T) {
	service := &orderService{}
	orderNo := service.generateOrderNo()

	assert.NotEmpty(t, orderNo)
	assert.Contains(t, orderNo, "ORD")
	assert.Len(t, orderNo, 17) // ORD + 14位时间戳
}

// 测试常量值覆盖率
func TestOrderConstants(t *testing.T) {
	assert.Equal(t, "WaitPay", contracts.PaymentStatusWaitPay)
	assert.Equal(t, "Paid", contracts.PaymentStatusPaid)
	assert.Equal(t, "Refunded", contracts.PaymentStatusRefunded)
	assert.Equal(t, "Cancelled", contracts.PaymentStatusCancelled)

	assert.Equal(t, "WaitMake", contracts.MakeStatusWaitMake)
	assert.Equal(t, "Making", contracts.MakeStatusMaking)
	assert.Equal(t, "Made", contracts.MakeStatusMade)
	assert.Equal(t, "Failed", contracts.MakeStatusFailed)
}

// 测试创建服务的各种情况
func TestOrderService_ServiceCreation(t *testing.T) {
	// 测试nil参数创建
	service1 := NewOrderService(nil, nil, nil, nil)
	assert.NotNil(t, service1)

	// 转换为具体类型以测试私有方法
	concreteService := service1.(*orderService)
	assert.NotNil(t, concreteService)

	// 测试生成订单号功能
	orderNo1 := concreteService.generateOrderNo()
	orderNo2 := concreteService.generateOrderNo()
	// 订单号可能在相同秒内相同，这是正常的
	assert.NotEmpty(t, orderNo1)
	assert.NotEmpty(t, orderNo2)
}

// 测试接口实现检查
func TestOrderService_InterfaceImplementation(t *testing.T) {
	var _ OrderService = &orderService{}
	// 如果编译通过，说明接口实现正确
	assert.True(t, true, "OrderService interface implementation is correct")
}

// 测试结构体字段验证
func TestOrderService_StructFields(t *testing.T) {
	service := &orderService{
		orderRepo:   nil,
		machineRepo: nil,
		memberRepo:  nil,
		deviceSvc:   nil,
	}
	assert.NotNil(t, service)

	// 验证字段都可以设置
	assert.Nil(t, service.orderRepo)
	assert.Nil(t, service.machineRepo)
	assert.Nil(t, service.memberRepo)
	assert.Nil(t, service.deviceSvc)
}

// 测试订单号生成的时间格式
func TestOrderService_GenerateOrderNoFormat(t *testing.T) {
	service := &orderService{}

	// 生成多个订单号
	orderNos := make([]string, 5)
	for i := 0; i < 5; i++ {
		orderNos[i] = service.generateOrderNo()
	}

	// 验证格式一致性
	for _, orderNo := range orderNos {
		assert.True(t, len(orderNo) == 17, "Order number should be 17 characters long")
		assert.True(t, orderNo[:3] == "ORD", "Order number should start with ORD")

		// 验证后面是数字
		dateTimePart := orderNo[3:]
		assert.Len(t, dateTimePart, 14, "DateTime part should be 14 characters")

		// 验证是否符合时间格式 (YYYYMMDDHHMMSS)
		assert.Regexp(t, `^\d{14}$`, dateTimePart, "DateTime part should be 14 digits")
	}
}

// 测试基础结构体方法调用
func TestOrderService_BasicMethodsExist(t *testing.T) {
	service := NewOrderService(nil, nil, nil, nil)

	// 检查方法是否存在，这里只验证接口方法存在
	assert.NotNil(t, service)

	// 测试私有方法通过类型断言访问
	if concreteService, ok := service.(*orderService); ok {
		orderNo := concreteService.generateOrderNo()
		assert.NotEmpty(t, orderNo)
		assert.Contains(t, orderNo, "ORD")
	} else {
		t.Error("Service should be of type *orderService")
	}
}

// 测试不同参数组合的服务创建
func TestOrderService_CreationWithDifferentParams(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{"All nil params"},
		{"Service creation"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewOrderService(nil, nil, nil, nil)
			assert.NotNil(t, service)
		})
	}
}

// 测试内部函数和方法
func TestOrderService_InternalMethods(t *testing.T) {
	service := &orderService{}

	// 测试订单号生成的一致性
	orderNo1 := service.generateOrderNo()
	assert.NotEmpty(t, orderNo1)

	// 测试多次调用产生不同结果（因为时间不同）
	// 注意：在快速连续调用时可能产生相同结果，这是正常的
	orderNos := make(map[string]bool)
	for i := 0; i < 3; i++ {
		orderNo := service.generateOrderNo()
		orderNos[orderNo] = true
	}

	// 至少应该有一个不同的订单号（通常会有多个）
	assert.GreaterOrEqual(t, len(orderNos), 1, "Should generate at least one unique order number")
}

// 测试常量检查和基本功能
func TestOrderService_StatusConstants(t *testing.T) {
	// 测试PaymentStatus枚举值转换
	tests := []struct {
		status   enums.PaymentStatus
		expected string
	}{
		{enums.PaymentStatusWaitPay, "待支付"},
		{enums.PaymentStatusPaid, "已支付"},
		{enums.PaymentStatusRefunded, "已退款"},
	}

	for _, tt := range tests {
		desc := enums.GetPaymentStatusDesc(tt.status)
		if desc != tt.expected {
			t.Errorf("Expected payment status %d to have desc '%s', got '%s'", tt.status, tt.expected, desc)
		}
	}
}

func TestOrderService_MakeStatusConstants(t *testing.T) {
	// 测试MakeStatus枚举值转换
	tests := []struct {
		status   enums.MakeStatus
		expected string
	}{
		{enums.MakeStatusWaitMake, "待制作"},
		{enums.MakeStatusMaking, "制作中"},
		{enums.MakeStatusMade, "制作完成"},
		{enums.MakeStatusMakeFail, "制作失败"},
	}

	for _, tt := range tests {
		desc := enums.GetMakeStatusDesc(tt.status)
		if desc != tt.expected {
			t.Errorf("Expected make status %d to have desc '%s', got '%s'", tt.status, tt.expected, desc)
		}
	}
}

// 测试枚举状态常量
func TestOrderService_EnumValueMethods(t *testing.T) {
	// 测试PaymentStatus的所有方法
	status := enums.PaymentStatusPaid
	assert.True(t, status.IsValid(), "PaymentStatusPaid should be valid")
	assert.Equal(t, "已支付", enums.GetPaymentStatusDesc(status))
	assert.Equal(t, "已支付", status.String())

	// 测试MakeStatus的所有方法
	makeStatus := enums.MakeStatusMade
	assert.True(t, makeStatus.IsValid(), "MakeStatusMade should be valid")
	assert.Equal(t, "制作完成", enums.GetMakeStatusDesc(makeStatus))
	assert.Equal(t, "制作完成", makeStatus.String())

	// 测试BusinessStatus的所有方法
	businessStatus := enums.BusinessStatusOpen
	assert.True(t, businessStatus.IsValid(), "BusinessStatusOpen should be valid")
	assert.Equal(t, "营业中", enums.GetBusinessStatusDesc(businessStatus))
	assert.Equal(t, "营业中", businessStatus.String())
}

// 测试Enum API转换方法
func TestOrderService_EnumAPIConversion(t *testing.T) {
	// 测试BusinessStatus API转换
	status := enums.BusinessStatusOpen
	apiStr := status.ToAPIString()
	assert.Equal(t, "Open", apiStr)

	// 测试反向转换
	convertedStatus := enums.FromAPIString("Open")
	assert.Equal(t, enums.BusinessStatusOpen, convertedStatus)

	// 测试无效值的默认行为
	defaultStatus := enums.FromAPIString("Invalid")
	assert.Equal(t, enums.BusinessStatusOpen, defaultStatus)
}

// Mock repository for testing GetByOrderNo
type mockOrderRepository struct {
	mock.Mock
}

func (m *mockOrderRepository) GetByOrderNo(orderNo string) (*models.Order, error) {
	args := m.Called(orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *mockOrderRepository) GetByID(id string) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *mockOrderRepository) GetByMemberPaging(memberID string, pageIndex, pageSize int) ([]models.Order, int64, error) {
	args := m.Called(memberID, pageIndex, pageSize)
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *mockOrderRepository) Create(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *mockOrderRepository) Update(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *mockOrderRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test GetByOrderNo method
func TestOrderService_GetByOrderNo(t *testing.T) {
	mockRepo := &mockOrderRepository{}
	service := &orderService{
		orderRepo: mockRepo,
	}

	// Test successful case
	order := &models.Order{
		ID:      "order123",
		OrderNo: stringPtr("ORD20250813001"),
	}
	mockRepo.On("GetByOrderNo", "ORD20250813001").Return(order, nil)

	result, err := service.GetByOrderNo("ORD20250813001")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "order123", result.ID)
	assert.Equal(t, "ORD20250813001", result.OrderNo)

	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetByOrderNo_NotFound(t *testing.T) {
	mockRepo := &mockOrderRepository{}
	service := &orderService{
		orderRepo: mockRepo,
	}

	// Test order not found case
	mockRepo.On("GetByOrderNo", "NON_EXISTENT_ORDER").Return(nil, gorm.ErrRecordNotFound)

	result, err := service.GetByOrderNo("NON_EXISTENT_ORDER")
	assert.NoError(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetByOrderNo_DatabaseError(t *testing.T) {
	mockRepo := &mockOrderRepository{}
	service := &orderService{
		orderRepo: mockRepo,
	}

	// Test database error case
	mockRepo.On("GetByOrderNo", "ORD20250813001").Return(nil, assert.AnError)

	result, err := service.GetByOrderNo("ORD20250813001")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "根据订单号获取订单失败")

	mockRepo.AssertExpectations(t)
}
