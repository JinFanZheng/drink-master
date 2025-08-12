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

func setupMemberTestRouter() (*gin.Engine, *MemberHandler) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	router := gin.New()
	memberHandler := NewMemberHandler(db)

	router.POST("/api/Member/Update", memberHandler.Update)
	router.POST("/api/Member/AddFranchiseIntention", memberHandler.AddFranchiseIntention)
	router.GET("/api/Member/GetUserInfo", memberHandler.GetUserInfo)

	return router, memberHandler
}

func TestNewMemberHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewMemberHandler(db)

	if handler == nil {
		t.Error("Expected handler to be created")
	}
}

func TestMemberHandler_Update(t *testing.T) {
	router, _ := setupMemberTestRouter()

	// 测试没有认证的情况
	updateData := map[string]interface{}{
		"nickname": "新用户名",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("POST", "/api/Member/Update", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMemberHandler_Update_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewMemberHandler(db)

	// 创建有效的JSON请求体
	updateData := map[string]interface{}{
		"nickname": "新用户名",
		"avatar":   "https://example.com/avatar.jpg",
	}
	jsonData, _ := json.Marshal(updateData)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Member/Update", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置认证信息
	c.Set("member_id", "test_member_123")

	handler.Update(c)

	// 由于数据库中没有这个member，期望会返回500错误
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestMemberHandler_AddFranchiseIntention(t *testing.T) {
	router, _ := setupMemberTestRouter()

	// 测试没有认证的情况
	intentionData := map[string]interface{}{
		"location":   "上海市",
		"investment": 100000,
	}

	jsonData, _ := json.Marshal(intentionData)
	req, _ := http.NewRequest("POST", "/api/Member/AddFranchiseIntention", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMemberHandler_AddFranchiseIntention_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewMemberHandler(db)

	// 创建有效的JSON请求体
	intentionData := map[string]interface{}{
		"contactName":      "张三",
		"contactPhone":     "13800138000",
		"intendedLocation": "北京市朝阳区",
		"remarks":          "希望开设新店",
	}
	jsonData, _ := json.Marshal(intentionData)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Member/AddFranchiseIntention", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// 设置认证信息
	c.Set("member_id", "test_member_456")

	handler.AddFranchiseIntention(c)

	// 由于数据库中没有这个member，期望会返回500错误
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestMemberHandler_GetUserInfo(t *testing.T) {
	router, _ := setupMemberTestRouter()

	// 测试没有认证的情况
	req, _ := http.NewRequest("GET", "/api/Member/GetUserInfo", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// 没有member_id会返回401
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMemberHandler_GetUserInfo_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewMemberHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/Member/GetUserInfo", nil)

	// 设置认证信息
	c.Set("member_id", "test_member_789")

	handler.GetUserInfo(c)

	// 由于数据库中没有这个member，期望会返回500错误
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
