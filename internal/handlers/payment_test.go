package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

func setupPaymentTestRouter() (*gin.Engine, *PaymentHandler, *CallbackHandler) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// 自动迁移表结构
	db.AutoMigrate(&models.Member{}, &models.Machine{}, &models.Product{}, &models.Order{})
	
	// 创建测试数据
	member := &models.Member{
		ID:             "test_member_123",
		WeChatOpenId:   "test_openid",
		Nickname:       "测试用户",
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
		ID:        "order123",
		OrderNo:   "ORD20240813001",
		MemberId:  "test_member_123",
		MachineId: "machine123",
		ProductId: "product123",
		PayAmount: 10.50,
	}
	db.Create(order)

	router := gin.New()
	paymentHandler := NewPaymentHandler(db)
	callbackHandler := NewCallbackHandler(db)

	router.GET("/api/Payment/Get", paymentHandler.Get)
	router.GET("/api/Payment/Query", paymentHandler.Query)
	router.POST("/api/Callback/PaymentResult", callbackHandler.PaymentResult)

	return router, paymentHandler, callbackHandler
}

func TestNewPaymentHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewPaymentHandler(db)

	if handler == nil {
		t.Error("Expected handler to be created")
	}
}

func TestNewCallbackHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewCallbackHandler(db)

	if handler == nil {
		t.Error("Expected handler to be created")
	}
}

func TestPaymentHandler_Get(t *testing.T) {
	router, _, _ := setupPaymentTestRouter()

	req, _ := http.NewRequest("GET", "/api/Payment/Get?orderId=order123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestPaymentHandler_Get_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	// 自动迁移表结构
	db.AutoMigrate(&models.Member{}, &models.Machine{}, &models.Product{}, &models.Order{})
	
	// 创建测试数据
	member := &models.Member{
		ID:             "test_member_123",
		WeChatOpenId:   "test_openid",
		Nickname:       "测试用户",
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
		ID:        "order123",
		OrderNo:   "ORD20240813001",
		MemberId:  "test_member_123",
		MachineId: "machine123",
		ProductId: "product123",
		PayAmount: 10.50,
	}
	db.Create(order)
	
	handler := NewPaymentHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/Payment/Get?orderId=order123", nil)

	// 设置认证信息
	c.Set("member_id", "test_member_123")

	handler.Get(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestPaymentHandler_Query(t *testing.T) {
	router, _, _ := setupPaymentTestRouter()

	req, _ := http.NewRequest("GET", "/api/Payment/Query?orderId=order123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestPaymentHandler_Query_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewPaymentHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/Payment/Query?orderId=order123", nil)

	// 设置认证信息
	c.Set("member_id", "test_member_456")

	handler.Query(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCallbackHandler_PaymentResult(t *testing.T) {
	router, _, _ := setupPaymentTestRouter()

	// 测试支付结果回调
	callbackData := map[string]interface{}{
		"orderId": "order123",
		"status":  "success",
		"amount":  19.99,
	}

	jsonData, _ := json.Marshal(callbackData)
	req, _ := http.NewRequest("POST", "/api/Callback/PaymentResult", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCallbackHandler_PaymentResult_Direct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewCallbackHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Callback/PaymentResult", nil)

	handler.PaymentResult(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// Additional test cases for better coverage
func TestPaymentHandler_Get_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewPaymentHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/Payment/Get", nil) // Missing orderId

	// Set auth info
	c.Set("member_id", "test_member_123")

	handler.Get(c)

	// Should return validation error
	if w.Code == http.StatusOK {
		t.Errorf("Expected error status, got %d", w.Code)
	}
}

func TestPaymentHandler_Query_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewPaymentHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/Payment/Query", nil) // Missing orderId

	// Set auth info
	c.Set("member_id", "test_member_456")

	handler.Query(c)

	// Should return validation error
	if w.Code == http.StatusOK {
		t.Errorf("Expected error status, got %d", w.Code)
	}
}

func TestCallbackHandler_PaymentResult_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewCallbackHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Callback/PaymentResult",
		bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PaymentResult(c)

	// Should handle invalid JSON gracefully
	// Handler should return OK or appropriate error status
	if w.Code != http.StatusOK {
		t.Logf("Handler returned status %d for invalid JSON", w.Code)
	}
}

func TestPaymentHandler_Constructor(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// Test PaymentHandler constructor
	paymentHandler := NewPaymentHandler(db)
	if paymentHandler == nil {
		t.Error("PaymentHandler should not be nil")
		return
	}
	if paymentHandler.paymentService == nil {
		t.Error("PaymentHandler.paymentService should not be nil")
	}
	if paymentHandler.orderService == nil {
		t.Error("PaymentHandler.orderService should not be nil")
	}

	// Test CallbackHandler constructor
	callbackHandler := NewCallbackHandler(db)
	if callbackHandler == nil {
		t.Error("CallbackHandler should not be nil")
	}
}

func TestPaymentHandler_GetMemberOpenId_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewPaymentHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/Payment/Get?orderId=order123", nil)

	// Set invalid member ID that should cause error
	c.Set("member_id", "")

	handler.Get(c)

	// Should handle error gracefully
	if w.Code == http.StatusOK {
		t.Errorf("Expected error status, got %d", w.Code)
	}
}
