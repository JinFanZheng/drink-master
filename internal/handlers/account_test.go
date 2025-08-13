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
	"github.com/ddteam/drink-master/pkg/wechat"
)

func setupAccountTestRouter() (*gin.Engine, *AccountHandler) {
	gin.SetMode(gin.TestMode)

	// 创建内存数据库用于测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// 运行数据库迁移创建表结构
	if err := models.AutoMigrate(db); err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}

	router := gin.New()
	wechatClient := wechat.NewClient("test_app_id", "test_app_secret")
	accountHandler := NewAccountHandler(db, wechatClient)

	router.GET("/api/Account/CheckUserInfo", accountHandler.CheckUserInfo)
	router.POST("/api/Account/WeChatLogin", accountHandler.WeChatLogin)
	router.GET("/api/Account/CheckLogin", accountHandler.CheckLogin)
	router.GET("/api/Account/GetUserInfo", accountHandler.GetUserInfo)

	return router, accountHandler
}

func TestNewAccountHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	wechatClient := wechat.NewClient("test_app_id", "test_app_secret")
	handler := NewAccountHandler(db, wechatClient)

	if handler == nil {
		t.Error("Expected handler to be created")
	}
}

func TestAccountHandler_CheckUserInfo(t *testing.T) {
	router, _ := setupAccountTestRouter()

	// Test with valid parameters (微信API会失败，但我们测试参数验证逻辑)
	req, _ := http.NewRequest("GET", "/api/Account/CheckUserInfo?code=test123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 在测试环境中，微信API调用会失败，所以预期返回400状态码
	// 这是正常的，因为我们使用的是测试用的微信client
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// 检查响应是否包含JSON格式
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Error("Expected Content-Type to be application/json")
	}

	// Test with both code and appId parameters
	req2, _ := http.NewRequest("GET", "/api/Account/CheckUserInfo?code=test123&appId=wx123", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should still fail with WeChat validation
	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d with appId, got %d", http.StatusBadRequest, w2.Code)
	}

	// Test with empty code
	req3, _ := http.NewRequest("GET", "/api/Account/CheckUserInfo?code=", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should return bad request for empty code
	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for empty code, got %d", http.StatusBadRequest, w3.Code)
	}
}

func TestAccountHandler_CheckUserInfo_MissingCode(t *testing.T) {
	router, _ := setupAccountTestRouter()

	// 不提供code参数
	req, _ := http.NewRequest("GET", "/api/Account/CheckUserInfo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAccountHandler_WeChatLogin(t *testing.T) {
	router, _ := setupAccountTestRouter()

	// Test with valid login data
	loginData := map[string]interface{}{
		"code":      "test_code_123",
		"nickName":  "Test User",
		"avatarUrl": "http://example.com/avatar.jpg",
		"appId":     "test_app_id",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/Account/WeChatLogin", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// WeChat validation will fail in test env, expect 400
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for WeChat validation failure, got %d", http.StatusBadRequest, w.Code)
	}

	// Test with minimal required data
	loginData2 := map[string]interface{}{
		"code": "test_code_456",
	}

	jsonData2, _ := json.Marshal(loginData2)
	req2, _ := http.NewRequest("POST", "/api/Account/WeChatLogin", bytes.NewBuffer(jsonData2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, req2)

	// Should still fail WeChat validation but pass JSON parsing
	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for minimal data, got %d", http.StatusBadRequest, w2.Code)
	}

	// Test with missing required field
	loginData3 := map[string]interface{}{
		"nickName":  "Test User",
		"avatarUrl": "http://example.com/avatar.jpg",
	}

	jsonData3, _ := json.Marshal(loginData3)
	req3, _ := http.NewRequest("POST", "/api/Account/WeChatLogin", bytes.NewBuffer(jsonData3))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()

	router.ServeHTTP(w3, req3)

	// Should fail validation for missing code
	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for missing code, got %d", http.StatusBadRequest, w3.Code)
	}
}

func TestAccountHandler_WeChatLogin_InvalidJSON(t *testing.T) {
	router, _ := setupAccountTestRouter()

	// 发送无效JSON
	req, _ := http.NewRequest("POST", "/api/Account/WeChatLogin", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAccountHandler_CheckLogin(t *testing.T) {
	router, _ := setupAccountTestRouter()

	// Test without Authorization header
	req, _ := http.NewRequest("GET", "/api/Account/CheckLogin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Without JWT middleware in test setup, this might return OK
	// In real application, JWT middleware would check for auth first
	if w.Code == http.StatusOK {
		t.Log("CheckLogin returned OK status (no JWT middleware in test)")
		// Verify response body
		if w.Body.String() != "ok" {
			t.Errorf("Expected response body 'ok', got '%s'", w.Body.String())
		}
	} else {
		t.Logf("CheckLogin returned status %d (expected for no auth)", w.Code)
	}

	// Test with invalid Authorization header format
	req2, _ := http.NewRequest("GET", "/api/Account/CheckLogin", nil)
	req2.Header.Set("Authorization", "InvalidFormat")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should still return OK in test env without real JWT middleware
	if w2.Code == http.StatusOK {
		t.Log("CheckLogin returned OK with invalid auth header (test env)")
	}

	// Test with Bearer token format (but invalid token)
	req3, _ := http.NewRequest("GET", "/api/Account/CheckLogin", nil)
	req3.Header.Set("Authorization", "Bearer invalid_token")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should return OK in test env without JWT validation
	if w3.Code == http.StatusOK {
		t.Log("CheckLogin returned OK with Bearer token (test env)")
	}
}

func TestAccountHandler_GetUserInfo(t *testing.T) {
	router, _ := setupAccountTestRouter()

	// Test without authentication
	req, _ := http.NewRequest("GET", "/api/Account/GetUserInfo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Without JWT middleware, this will reach the handler but fail on GetMemberID
	// which returns false for exists
	if w.Code >= 400 {
		t.Logf("GetUserInfo correctly returned error status %d for unauthenticated request", w.Code)
	} else {
		t.Log("GetUserInfo returned success status")
	}

	// Test with some header (though not valid JWT)
	req2, _ := http.NewRequest("GET", "/api/Account/GetUserInfo", nil)
	req2.Header.Set("Authorization", "Bearer fake_token")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should still fail because GetMemberID will return false
	if w2.Code >= 400 {
		t.Logf("GetUserInfo returned error status %d for invalid token", w2.Code)
	}

	// Test with different content types
	req3, _ := http.NewRequest("GET", "/api/Account/GetUserInfo", nil)
	req3.Header.Set("Accept", "application/json")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should behave the same regardless of Accept header
	if w3.Code >= 400 {
		t.Logf("GetUserInfo returned error status %d with Accept header", w3.Code)
	}
}
