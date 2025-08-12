package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// 创建内存数据库用于测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	router := gin.New()
	healthHandler := NewHealthHandler(db)

	router.GET("/health", healthHandler.Health)
	router.GET("/health/db", healthHandler.DatabaseHealth)

	return router, db
}

func TestNewHealthHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewHealthHandler(db)

	if handler.db != db {
		t.Error("Expected handler.db to be set correctly")
	}
}

func TestHealthHandler_Health(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
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

func TestHealthHandler_DatabaseHealth(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/health/db", nil)
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
