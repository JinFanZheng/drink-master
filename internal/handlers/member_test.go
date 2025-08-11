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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Member/Update", nil)

	// 设置认证信息
	c.Set("member_id", "test_member_123")

	handler.Update(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/Member/AddFranchiseIntention", nil)

	// 设置认证信息
	c.Set("member_id", "test_member_456")

	handler.AddFranchiseIntention(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
