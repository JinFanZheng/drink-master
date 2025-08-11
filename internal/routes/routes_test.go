package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}
	return db
}

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	
	router := SetupRoutes(db)
	
	if router == nil {
		t.Error("Expected router to be created")
	}
	
	// 测试健康检查路由
	req, _ := http.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSetupRoutes_DatabaseHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	
	router := SetupRoutes(db)
	
	// 测试数据库健康检查路由
	req, _ := http.NewRequest("GET", "/api/health/db", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}