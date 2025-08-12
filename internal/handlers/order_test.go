package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
)

// Mock OrderService for testing
type mockOrderService struct {
	mock.Mock
}

func (m *mockOrderService) GetMemberOrderPaging(request contracts.GetMemberOrderPagingRequest) (*contracts.OrderPagingResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contracts.OrderPagingResponse), args.Error(1)
}

func (m *mockOrderService) GetByID(id string) (*contracts.GetOrderByIdResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contracts.GetOrderByIdResponse), args.Error(1)
}

func (m *mockOrderService) Create(request contracts.CreateOrderRequest) (*contracts.CreateOrderResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contracts.CreateOrderResponse), args.Error(1)
}

func (m *mockOrderService) Refund(request contracts.RefundOrderRequest) (*contracts.RefundOrderResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*contracts.RefundOrderResponse), args.Error(1)
}

func setupOrderTestRouter() (*gin.Engine, *OrderHandler) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	router := gin.New()
	// 使用nil服务进行基本测试，在需要服务的测试中使用mock
	orderHandler := NewOrderHandler(db, nil)

	router.POST("/api/Order/GetPaging", orderHandler.GetPaging)
	router.GET("/api/Order/Get", orderHandler.Get)
	router.POST("/api/Order/Create", orderHandler.Create)
	router.POST("/api/Order/Refund", orderHandler.Refund)

	return router, orderHandler
}

func TestNewOrderHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	mockService := &mockOrderService{}
	handler := NewOrderHandler(db, mockService)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.orderService)
}

func TestOrderHandler_GetPaging(t *testing.T) {
	router, _ := setupOrderTestRouter()

	// 测试没有认证的情况
	pagingData := map[string]interface{}{
		"pageIndex": 1,
		"pageSize":  10,
	}

	jsonData, _ := json.Marshal(pagingData)
	req, _ := http.NewRequest("POST", "/api/Order/GetPaging", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOrderHandler_GetPaging_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// 创建mock服务
	mockService := &mockOrderService{}
	handler := NewOrderHandler(db, mockService)

	// 设置期望
	expectedResponse := &contracts.OrderPagingResponse{
		Orders: []contracts.GetMemberOrderPagingResponse{
			{
				ID:            "order-123",
				OrderNo:       "ORD202508120001",
				ProductName:   "拿铁咖啡",
				PayAmount:     decimal.NewFromFloat(15.80),
				CreatedAt:     time.Now(),
				PaymentStatus: "Paid",
			},
		},
		Meta: contracts.PaginationMeta{
			Total: 1,
			Count: 1,
			Meta: &contracts.Meta{
				Timestamp: time.Now(),
			},
		},
	}

	mockService.On("GetMemberOrderPaging", mock.AnythingOfType("contracts.GetMemberOrderPagingRequest")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 构建请求体
	requestBody := `{"pageIndex": 1, "pageSize": 10}`
	c.Request, _ = http.NewRequest("POST", "/api/Order/GetPaging", bytes.NewBuffer([]byte(requestBody)))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置认证信息
	c.Set("member_id", "test_member_123")

	handler.GetPaging(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestOrderHandler_Get(t *testing.T) {
	router, _ := setupOrderTestRouter()

	req, _ := http.NewRequest("GET", "/api/Order/Get?id=order123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOrderHandler_Get_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// 创建mock服务
	mockService := &mockOrderService{}
	handler := NewOrderHandler(db, mockService)

	// 设置期望
	expectedResponse := &contracts.GetOrderByIdResponse{
		ID:            "order-123",
		OrderNo:       "ORD202508120001",
		MachineID:     "machine-001",
		ProductName:   "拿铁咖啡",
		PayAmount:     decimal.NewFromFloat(15.80),
		PaymentStatus: "Paid",
		CreatedAt:     time.Now(),
	}

	mockService.On("GetByID", "order123").Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/Order/Get?id=order123", nil)

	// 设置认证信息
	c.Set("member_id", "test_member_456")

	handler.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestOrderHandler_Create(t *testing.T) {
	router, _ := setupOrderTestRouter()

	// 测试没有认证的情况
	orderData := map[string]interface{}{
		"machineId": "machine123",
		"productId": "product456",
		"hasCup":    true,
		"payAmount": "15.80",
	}

	jsonData, _ := json.Marshal(orderData)
	req, _ := http.NewRequest("POST", "/api/Order/Create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestOrderHandler_Create_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// 创建mock服务
	mockService := &mockOrderService{}
	handler := NewOrderHandler(db, mockService)

	// 设置期望
	expectedResponse := &contracts.CreateOrderResponse{
		OrderID: "order-123",
		OrderNo: "ORD202508120001",
		Message: "订单创建成功",
	}

	mockService.On("Create", mock.AnythingOfType("contracts.CreateOrderRequest")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 构建请求体
	requestBody := `{"machineId": "machine123", "productId": "product456", "hasCup": true, "payAmount": "15.80"}`
	c.Request, _ = http.NewRequest("POST", "/api/Order/Create", bytes.NewBuffer([]byte(requestBody)))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置认证信息
	c.Set("member_id", "test_member_789")

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestOrderHandler_Refund(t *testing.T) {
	router, _ := setupOrderTestRouter()

	// 测试没有认证的情况
	refundData := map[string]interface{}{
		"orderId": "order123",
		"reason":  "商品有问题",
	}

	jsonData, _ := json.Marshal(refundData)
	req, _ := http.NewRequest("POST", "/api/Order/Refund", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回403，因为退款需要机主权限
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestOrderHandler_Refund_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// 创建mock服务
	mockService := &mockOrderService{}
	handler := NewOrderHandler(db, mockService)

	// 设置期望
	expectedResponse := &contracts.RefundOrderResponse{
		OrderID:      "order-123",
		RefundAmount: decimal.NewFromFloat(15.80),
		Message:      "退款成功",
	}

	mockService.On("Refund", mock.AnythingOfType("contracts.RefundOrderRequest")).Return(expectedResponse, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 构建请求体
	requestBody := `{"orderId": "order123", "reason": "商品有问题"}`
	c.Request, _ = http.NewRequest("POST", "/api/Order/Refund", bytes.NewBuffer([]byte(requestBody)))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置认证信息和机主权限
	c.Set("member_id", "test_member_101")
	c.Set("role", "Owner")

	handler.Refund(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
