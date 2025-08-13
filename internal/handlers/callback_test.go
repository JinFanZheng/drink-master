package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
	"github.com/ddteam/drink-master/internal/repositories"
	"github.com/ddteam/drink-master/internal/services"
)

// Mock services for testing
type mockCallbackOrderService struct {
	mock.Mock
}

func (m *mockCallbackOrderService) GetMemberOrderPaging(request contracts.GetMemberOrderPagingRequest) (*contracts.OrderPagingResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*contracts.OrderPagingResponse), args.Error(1)
}

func (m *mockCallbackOrderService) GetByID(id string) (*contracts.GetOrderByIdResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*contracts.GetOrderByIdResponse), args.Error(1)
}

func (m *mockCallbackOrderService) GetByOrderNo(orderNo string) (*models.Order, error) {
	args := m.Called(orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *mockCallbackOrderService) Create(request contracts.CreateOrderRequest) (*contracts.CreateOrderResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*contracts.CreateOrderResponse), args.Error(1)
}

func (m *mockCallbackOrderService) Refund(request contracts.RefundOrderRequest) (*contracts.RefundOrderResponse, error) {
	args := m.Called(request)
	return args.Get(0).(*contracts.RefundOrderResponse), args.Error(1)
}

type mockPaymentService struct {
	mock.Mock
}

func (m *mockPaymentService) WeChatPay(req contracts.WeChatPayRequest) (*contracts.WeChatPayResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*contracts.WeChatPayResponse), args.Error(1)
}

func (m *mockPaymentService) TranQuery(req contracts.TranQueryRequest) (*contracts.TranQueryResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*contracts.TranQueryResponse), args.Error(1)
}

func (m *mockPaymentService) GetPaymentAccount(machineID string) (*contracts.PaymentAccount, error) {
	args := m.Called(machineID)
	return args.Get(0).(*contracts.PaymentAccount), args.Error(1)
}

func (m *mockPaymentService) PayOrder(req contracts.PayOrderRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *mockPaymentService) InvalidOrder(req contracts.InvalidOrderRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *mockPaymentService) ProcessPaymentCallback(req contracts.PaymentCallbackRequest) (*contracts.PaymentCallbackResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*contracts.PaymentCallbackResponse), args.Error(1)
}

func setupCallbackTestRouter() (*gin.Engine, *CallbackHandler) {
	gin.SetMode(gin.TestMode)

	mockOrderSvc := &mockCallbackOrderService{}
	mockPaymentSvc := &mockPaymentService{}
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel) // Reduce log noise in tests

	handler := NewCallbackHandler(mockOrderSvc, mockPaymentSvc, logger)

	router := gin.New()
	router.POST("/api/Callback/PaymentResult", handler.PaymentResult)

	return router, handler
}

func TestNewCallbackHandler(t *testing.T) {
	mockOrderSvc := &mockCallbackOrderService{}
	mockPaymentSvc := &mockPaymentService{}
	logger := logrus.New()

	handler := NewCallbackHandler(mockOrderSvc, mockPaymentSvc, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, mockOrderSvc, handler.orderService)
	assert.Equal(t, mockPaymentSvc, handler.paymentService)
	assert.Equal(t, logger, handler.logger)
}

func TestCallbackHandler_PaymentResult_Success(t *testing.T) {
	router, handler := setupCallbackTestRouter()

	// 准备测试数据
	request := contracts.PaymentCallbackResultRequest{
		ChannelCode:    "test_channel",
		TransAmt:       1000,
		ReturnAmt:      0,
		OrderNo:        "ORD20250813001",
		OrderInfo:      "test order",
		ModeOfPayment:  1,
		ChannelOrderNo: "CH123456789",
		PaymentTime:    time.Now(),
		CallbackType:   "payment_success",
	}

	// Mock order (waiting for payment)
	order := &models.Order{
		ID:            "order123",
		OrderNo:       "ORD20250813001",
		PaymentStatus: int(enums.PaymentStatusWaitPay),
	}

	// Setup mocks
	mockOrderSvc := handler.orderService.(*mockCallbackOrderService)
	mockPaymentSvc := handler.paymentService.(*mockPaymentService)

	mockOrderSvc.On("GetByOrderNo", "ORD20250813001").Return(order, nil)
	mockPaymentSvc.On("PayOrder", mock.AnythingOfType("contracts.PayOrderRequest")).Return(nil)

	// 准备请求
	jsonData, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/Callback/PaymentResult", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 验证结果
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())

	// 验证mock调用
	mockOrderSvc.AssertCalled(t, "GetByOrderNo", "ORD20250813001")
	mockPaymentSvc.AssertCalled(t, "PayOrder", mock.AnythingOfType("contracts.PayOrderRequest"))
}

func TestCallbackHandler_PaymentResult_InvalidJSON(t *testing.T) {
	router, _ := setupCallbackTestRouter()

	// 发送无效JSON
	req, _ := http.NewRequest("POST", "/api/Callback/PaymentResult", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "参数错误", w.Body.String())
}

func TestCallbackHandler_PaymentResult_OrderNotFound(t *testing.T) {
	router, handler := setupCallbackTestRouter()

	request := contracts.PaymentCallbackResultRequest{
		ChannelCode:    "test_channel",
		TransAmt:       1000,
		OrderNo:        "NON_EXISTENT_ORDER",
		ChannelOrderNo: "CH123456789",
		PaymentTime:    time.Now(),
		ModeOfPayment:  1,
		CallbackType:   "payment_success",
	}

	// Setup mocks
	mockOrderSvc := handler.orderService.(*mockCallbackOrderService)
	mockOrderSvc.On("GetByOrderNo", "NON_EXISTENT_ORDER").Return(nil, nil) // Order not found

	jsonData, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/Callback/PaymentResult", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "订单不存在", w.Body.String())
	mockOrderSvc.AssertCalled(t, "GetByOrderNo", "NON_EXISTENT_ORDER")
}

func TestCallbackHandler_PaymentResult_OrderAlreadyPaid(t *testing.T) {
	router, handler := setupCallbackTestRouter()

	request := contracts.PaymentCallbackResultRequest{
		ChannelCode:    "test_channel",
		TransAmt:       1000,
		OrderNo:        "ORD20250813001",
		ChannelOrderNo: "CH123456789",
		PaymentTime:    time.Now(),
		ModeOfPayment:  1,
		CallbackType:   "payment_success",
	}

	// Mock order (already paid)
	order := &models.Order{
		ID:            "order123",
		OrderNo:       "ORD20250813001",
		PaymentStatus: int(enums.PaymentStatusPaid), // Already paid
	}

	// Setup mocks
	mockOrderSvc := handler.orderService.(*mockCallbackOrderService)
	mockOrderSvc.On("GetByOrderNo", "ORD20250813001").Return(order, nil)

	jsonData, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/Callback/PaymentResult", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
	mockOrderSvc.AssertCalled(t, "GetByOrderNo", "ORD20250813001")
}

func TestCallbackHandler_PaymentResult_PaymentServiceError(t *testing.T) {
	router, handler := setupCallbackTestRouter()

	request := contracts.PaymentCallbackResultRequest{
		ChannelCode:    "test_channel",
		TransAmt:       1000,
		OrderNo:        "ORD20250813001",
		ChannelOrderNo: "CH123456789",
		PaymentTime:    time.Now(),
		ModeOfPayment:  1,
		CallbackType:   "payment_success",
	}

	// Mock order (waiting for payment)
	order := &models.Order{
		ID:            "order123",
		OrderNo:       "ORD20250813001",
		PaymentStatus: int(enums.PaymentStatusWaitPay),
	}

	// Setup mocks
	mockOrderSvc := handler.orderService.(*mockCallbackOrderService)
	mockPaymentSvc := handler.paymentService.(*mockPaymentService)

	mockOrderSvc.On("GetByOrderNo", "ORD20250813001").Return(order, nil)
	mockPaymentSvc.On("PayOrder", mock.AnythingOfType("contracts.PayOrderRequest")).Return(assert.AnError)

	jsonData, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/Callback/PaymentResult", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Even if PayOrder fails, should return "ok" to prevent retry
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
	mockOrderSvc.AssertCalled(t, "GetByOrderNo", "ORD20250813001")
	mockPaymentSvc.AssertCalled(t, "PayOrder", mock.AnythingOfType("contracts.PayOrderRequest"))
}

// Integration test with real database
func TestCallbackHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto migrate
	err = db.AutoMigrate(&models.Member{}, &models.Machine{}, &models.Product{}, &models.Order{})
	assert.NoError(t, err)

	// Create test data
	member := &models.Member{
		ID:           "test_member_123",
		WeChatOpenId: "test_openid",
		Nickname:     "测试用户",
	}
	db.Create(member)

	machine := &models.Machine{
		ID:        "machine123",
		MachineNo: "VM001",
		Name:      "测试咖啡机",
	}
	db.Create(machine)

	product := &models.Product{
		ID:   "product123",
		Name: "测试饮品",
	}
	db.Create(product)

	order := &models.Order{
		ID:            "order123",
		OrderNo:       "ORD20250813001",
		MemberId:      "test_member_123",
		MachineId:     "machine123",
		ProductId:     "product123",
		PayAmount:     10.50,
		PaymentStatus: int(enums.PaymentStatusWaitPay),
	}
	db.Create(order)

	// Setup real services
	orderRepo := repositories.NewOrderRepository(db)
	machineRepo := repositories.NewMachineRepository(db)
	memberRepo := repositories.NewMemberRepository(db)
	deviceSvc := services.NewDeviceService()
	orderService := services.NewOrderService(orderRepo, machineRepo, memberRepo, deviceSvc)

	paymentService := services.NewPaymentService(db)

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	handler := NewCallbackHandler(orderService, paymentService, logger)

	router := gin.New()
	router.POST("/api/Callback/PaymentResult", handler.PaymentResult)

	// Test successful payment callback
	request := contracts.PaymentCallbackResultRequest{
		ChannelCode:    "test_channel",
		TransAmt:       1050, // 10.50 * 100
		OrderNo:        "ORD20250813001",
		ChannelOrderNo: "CH123456789",
		PaymentTime:    time.Now(),
		ModeOfPayment:  1,
		CallbackType:   "payment_success",
	}

	jsonData, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/api/Callback/PaymentResult", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())

	// Verify order was updated
	var updatedOrder models.Order
	db.First(&updatedOrder, "id = ?", "order123")
	assert.Equal(t, int(enums.PaymentStatusPaid), updatedOrder.PaymentStatus)
	assert.NotNil(t, updatedOrder.PaymentTime)
	assert.Equal(t, "CH123456789", *updatedOrder.ChannelOrderNo)
}
