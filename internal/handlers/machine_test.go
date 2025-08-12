package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

func setupMachineTestRouter() (*gin.Engine, *MachineHandler) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// 运行数据库迁移创建表结构
	if err := models.AutoMigrate(db); err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}

	router := gin.New()
	machineHandler := NewMachineHandler(db)

	router.GET("/api/Machine/Get", machineHandler.Get)
	router.GET("/api/Machine/CheckDeviceExist", machineHandler.CheckDeviceExist)
	router.GET("/api/Machine/GetProductList", machineHandler.GetProductList)
	router.POST("/api/Machine/GetPaging", machineHandler.GetPaging)
	router.GET("/api/Machine/GetList", machineHandler.GetList)
	router.GET("/api/Machine/OpenOrClose", machineHandler.OpenOrCloseBusiness)

	return router, machineHandler
}

func TestNewMachineHandler(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	handler := NewMachineHandler(db)

	if handler == nil {
		t.Error("Expected handler to be created")
	}
}

func TestMachineHandler_Get(t *testing.T) {
	router, _ := setupMachineTestRouter()

	// Test with machineNo parameter
	req, _ := http.NewRequest("GET", "/api/Machine/Get?machineNo=test123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status OK, NotFound, or BadRequest, got %d", w.Code)
	}

	// Test without machineNo parameter
	req2, _ := http.NewRequest("GET", "/api/Machine/Get", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should return bad request due to missing machineNo
	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest for missing machineNo, got %d", w2.Code)
	}

	// Test with invalid machineNo (empty string)
	req3, _ := http.NewRequest("GET", "/api/Machine/Get?machineNo=", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should return bad request for empty machineNo
	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest for empty machineNo, got %d", w3.Code)
	}
}

func TestMachineHandler_CheckDeviceExist(t *testing.T) {
	router, _ := setupMachineTestRouter()

	req, _ := http.NewRequest("GET", "/api/Machine/CheckDeviceExist?deviceId=device123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status OK or BadRequest, got %d", w.Code)
	}
}

func TestMachineHandler_GetProductList(t *testing.T) {
	router, _ := setupMachineTestRouter()

	req, _ := http.NewRequest("GET", "/api/Machine/GetProductList?machineId=machine123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status OK or BadRequest, got %d", w.Code)
	}
}

func TestMachineHandler_GetPaging(t *testing.T) {
	router, _ := setupMachineTestRouter()

	// Test without authentication
	req, _ := http.NewRequest("POST", "/api/Machine/GetPaging", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return forbidden or unauthorized
	if w.Code != http.StatusForbidden && w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Forbidden or Unauthorized, got %d", w.Code)
	}
}

func TestMachineHandler_GetList(t *testing.T) {
	router, _ := setupMachineTestRouter()

	// Test without authentication
	req, _ := http.NewRequest("GET", "/api/Machine/GetList", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return forbidden or unauthorized
	if w.Code != http.StatusForbidden && w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Forbidden or Unauthorized, got %d", w.Code)
	}
}

func TestMachineHandler_OpenOrCloseBusiness(t *testing.T) {
	router, _ := setupMachineTestRouter()

	// Test without authentication but with parameters
	req, _ := http.NewRequest("GET", "/api/Machine/OpenOrClose?machineId=machine123&status=Open", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return forbidden or unauthorized due to lack of authentication
	if w.Code != http.StatusForbidden && w.Code != http.StatusUnauthorized && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status Forbidden, Unauthorized, or BadRequest, got %d", w.Code)
	}

	// Test without machineId parameter
	req2, _ := http.NewRequest("GET", "/api/Machine/OpenOrClose?status=Open", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should return bad request or authentication error
	if w2.Code != http.StatusBadRequest && w2.Code != http.StatusForbidden && w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status BadRequest, Forbidden, or Unauthorized, got %d", w2.Code)
	}

	// Test without status parameter
	req3, _ := http.NewRequest("GET", "/api/Machine/OpenOrClose?machineId=machine123", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should return bad request or authentication error
	if w3.Code != http.StatusBadRequest && w3.Code != http.StatusForbidden && w3.Code != http.StatusUnauthorized {
		t.Errorf("Expected status BadRequest, Forbidden, or Unauthorized, got %d", w3.Code)
	}
}
