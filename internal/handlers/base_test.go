package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupBaseTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/test", nil)
	return c, w
}

func TestNewBaseHandler(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	handler := NewBaseHandler(db)
	
	if handler.db != db {
		t.Error("Expected handler.db to be set correctly")
	}
}

func TestBaseHandler_GetMemberID(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewBaseHandler(db)
	c, _ := setupBaseTestContext()
	
	// 测试没有member_id的情况
	memberID, exists := handler.GetMemberID(c)
	if exists {
		t.Error("Expected GetMemberID to return false when no member_id is set")
	}
	if memberID != "" {
		t.Error("Expected memberID to be empty when not exists")
	}
	
	// 设置member_id并测试
	c.Set("member_id", "test_member_123")
	memberID, exists = handler.GetMemberID(c)
	if !exists {
		t.Error("Expected GetMemberID to return true when member_id is set")
	}
	if memberID != "test_member_123" {
		t.Errorf("Expected memberID to be 'test_member_123', got '%s'", memberID)
	}
}

func TestBaseHandler_GetMachineOwnerID(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewBaseHandler(db)
	c, _ := setupBaseTestContext()
	
	// 测试没有machine_owner_id的情况
	ownerID, exists := handler.GetMachineOwnerID(c)
	if exists {
		t.Error("Expected GetMachineOwnerID to return false when no machine_owner_id is set")
	}
	if ownerID != "" {
		t.Error("Expected ownerID to be empty when not exists")
	}
	
	// 设置machine_owner_id并测试
	c.Set("machine_owner_id", "test_owner_456")
	ownerID, exists = handler.GetMachineOwnerID(c)
	if !exists {
		t.Error("Expected GetMachineOwnerID to return true when machine_owner_id is set")
	}
	if ownerID != "test_owner_456" {
		t.Errorf("Expected ownerID to be 'test_owner_456', got '%s'", ownerID)
	}
}

func TestBaseHandler_IsMachineOwner(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewBaseHandler(db)
	c, _ := setupBaseTestContext()
	
	// 测试没有role的情况
	if handler.IsMachineOwner(c) {
		t.Error("Expected IsMachineOwner to return false when no role is set")
	}
	
	// 设置非Owner角色
	c.Set("role", "Member")
	if handler.IsMachineOwner(c) {
		t.Error("Expected IsMachineOwner to return false for Member role")
	}
	
	// 设置Owner角色
	c.Set("role", "Owner")
	if !handler.IsMachineOwner(c) {
		t.Error("Expected IsMachineOwner to return true for Owner role")
	}
}

func TestBaseHandler_GetCurrentRole(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewBaseHandler(db)
	c, _ := setupBaseTestContext()
	
	// 测试没有role的情况
	role, exists := handler.GetCurrentRole(c)
	if exists {
		t.Error("Expected GetCurrentRole to return false when no role is set")
	}
	if role != "" {
		t.Error("Expected role to be empty when not exists")
	}
	
	// 设置role并测试
	c.Set("role", "Member")
	role, exists = handler.GetCurrentRole(c)
	if !exists {
		t.Error("Expected GetCurrentRole to return true when role is set")
	}
	if role != "Member" {
		t.Errorf("Expected role to be 'Member', got '%s'", role)
	}
}

func TestBaseHandler_ValidationErrorResponse(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewBaseHandler(db)
	c, w := setupBaseTestContext()
	
	err := fmt.Errorf("validation failed")
	handler.ValidationErrorResponse(c, err)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestBaseHandler_NotFoundResponse(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewBaseHandler(db)
	c, w := setupBaseTestContext()
	
	handler.NotFoundResponse(c, "Resource not found")
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestBaseHandler_ForbiddenResponse(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewBaseHandler(db)
	c, w := setupBaseTestContext()
	
	handler.ForbiddenResponse(c, "Access denied")
	
	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestBaseHandler_InternalErrorResponse(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewBaseHandler(db)
	c, w := setupBaseTestContext()
	
	err := fmt.Errorf("internal error")
	handler.InternalErrorResponse(c, err)
	
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}