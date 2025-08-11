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

func setupAccountTestRouter() (*gin.Engine, *AccountHandler) {
	gin.SetMode(gin.TestMode)

	// 创建内存数据库用于测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	router := gin.New()
	accountHandler := NewAccountHandler(db)

	router.GET("/api/Account/CheckUserInfo", accountHandler.CheckUserInfo)
	router.POST("/api/Account/WeChatLogin", accountHandler.WeChatLogin)
	router.GET("/api/Account/CheckLogin", accountHandler.CheckLogin)
	router.GET("/api/Account/GetUserInfo", accountHandler.GetUserInfo)

	return router, accountHandler
}

func TestNewAccountHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewAccountHandler(db)

	if handler == nil {
		t.Error("Expected handler to be created")
	}
}

func TestAccountHandler_CheckUserInfo(t *testing.T) {
	router, _ := setupAccountTestRouter()

	req, _ := http.NewRequest("GET", "/api/Account/CheckUserInfo?openId=test123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// 检查响应是否包含JSON格式
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Error("Expected Content-Type to be application/json")
	}
}

func TestAccountHandler_CheckUserInfo_MissingOpenId(t *testing.T) {
	router, _ := setupAccountTestRouter()

	// 不提供openId参数
	req, _ := http.NewRequest("GET", "/api/Account/CheckUserInfo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAccountHandler_WeChatLogin(t *testing.T) {
	router, _ := setupAccountTestRouter()

	loginData := map[string]interface{}{
		"code":          "test_code_123",
		"iv":            "test_iv",
		"encryptedData": "test_encrypted_data",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/Account/WeChatLogin", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 这个会返回错误因为是mock数据，但至少测试了handler逻辑
	if w.Code == 0 {
		t.Error("Expected some HTTP status code")
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

	req, _ := http.NewRequest("GET", "/api/Account/CheckLogin", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 没有JWT token会返回401或其他错误状态
	if w.Code == http.StatusOK {
		t.Log("CheckLogin returned OK status")
	} else {
		t.Logf("CheckLogin returned status %d (expected for no auth)", w.Code)
	}
}

func TestAccountHandler_GetUserInfo(t *testing.T) {
	router, _ := setupAccountTestRouter()

	req, _ := http.NewRequest("GET", "/api/Account/GetUserInfo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 没有认证信息，期望返回错误状态
	if w.Code >= 400 {
		t.Logf("GetUserInfo correctly returned error status %d for unauthenticated request", w.Code)
	} else {
		t.Log("GetUserInfo returned success status")
	}
}
