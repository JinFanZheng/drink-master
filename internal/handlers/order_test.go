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
)

func setupOrderTestRouter() (*gin.Engine, *OrderHandler) {
	gin.SetMode(gin.TestMode)
	
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	router := gin.New()
	orderHandler := NewOrderHandler(db)
	
	router.POST("/api/Order/GetPaging", orderHandler.GetPaging)
	router.GET("/api/Order/Get", orderHandler.Get)
	router.POST("/api/Order/Create", orderHandler.Create)
	router.POST("/api/Order/Refund", orderHandler.Refund)
	
	return router, orderHandler
}

func TestNewOrderHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewOrderHandler(db)
	
	if handler == nil {
		t.Error("Expected handler to be created")
	}
}

func TestOrderHandler_GetPaging(t *testing.T) {
	router, _ := setupOrderTestRouter()
	
	// 测试没有认证的情况
	pagingData := map[string]interface{}{
		"pageIndex": 1,
		"pageSize": 10,
	}
	
	jsonData, _ := json.Marshal(pagingData)
	req, _ := http.NewRequest("POST", "/api/Order/GetPaging", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOrderHandler_GetPaging_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewOrderHandler(db)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Order/GetPaging", nil)
	
	// 设置认证信息
	c.Set("member_id", "test_member_123")
	
	handler.GetPaging(c)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOrderHandler_Get(t *testing.T) {
	router, _ := setupOrderTestRouter()
	
	req, _ := http.NewRequest("GET", "/api/Order/Get?orderId=order123", nil)
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOrderHandler_Get_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewOrderHandler(db)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/Order/Get?orderId=order123", nil)
	
	// 设置认证信息
	c.Set("member_id", "test_member_456")
	
	handler.Get(c)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOrderHandler_Create(t *testing.T) {
	router, _ := setupOrderTestRouter()
	
	// 测试没有认证的情况
	orderData := map[string]interface{}{
		"machineId": "machine123",
		"productId": "product456",
		"hasCup": true,
	}
	
	jsonData, _ := json.Marshal(orderData)
	req, _ := http.NewRequest("POST", "/api/Order/Create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOrderHandler_Create_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewOrderHandler(db)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Order/Create", nil)
	
	// 设置认证信息
	c.Set("member_id", "test_member_789")
	
	handler.Create(c)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestOrderHandler_Refund(t *testing.T) {
	router, _ := setupOrderTestRouter()
	
	// 测试没有认证的情况
	refundData := map[string]interface{}{
		"orderId": "order123",
		"reason": "商品有问题",
	}
	
	jsonData, _ := json.Marshal(refundData)
	req, _ := http.NewRequest("POST", "/api/Order/Refund", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	router.ServeHTTP(w, req)
	
	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestOrderHandler_Refund_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewOrderHandler(db)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Order/Refund", nil)
	
	// 设置认证信息
	c.Set("member_id", "test_member_101")
	
	handler.Refund(c)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}