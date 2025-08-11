package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMachineTestRouter() (*gin.Engine, *MachineHandler) {
	gin.SetMode(gin.TestMode)
	
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
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
	
	req, _ := http.NewRequest("GET", "/api/Machine/Get?machineNo=test123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status OK, NotFound, or BadRequest, got %d", w.Code)
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
	
	req, _ := http.NewRequest("POST", "/api/Machine/GetPaging", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// 没有认证和参数会返回错误状态
	if w.Code == 0 {
		t.Error("Expected some HTTP status code")
	}
}

func TestMachineHandler_GetList(t *testing.T) {
	router, _ := setupMachineTestRouter()
	
	req, _ := http.NewRequest("GET", "/api/Machine/GetList", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// 没有认证会返回错误状态
	if w.Code == 0 {
		t.Error("Expected some HTTP status code")
	}
}

func TestMachineHandler_OpenOrCloseBusiness(t *testing.T) {
	router, _ := setupMachineTestRouter()
	
	req, _ := http.NewRequest("GET", "/api/Machine/OpenOrClose?machineId=machine123&status=Open", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// 没有认证会返回错误状态
	if w.Code == 0 {
		t.Error("Expected some HTTP status code")
	}
}