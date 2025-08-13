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

	// Test with valid id parameter
	req, _ := http.NewRequest("GET", "/api/Machine/Get?id=test123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status OK, NotFound, or InternalServerError, got %d", w.Code)
	}

	// Test without id parameter
	req2, _ := http.NewRequest("GET", "/api/Machine/Get", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should return bad request due to missing id
	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest for missing id, got %d", w2.Code)
	}

	// Test with empty id parameter
	req3, _ := http.NewRequest("GET", "/api/Machine/Get?id=", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should return bad request for empty id
	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest for empty id, got %d", w3.Code)
	}

	// Test with nonexistent machine id
	req4, _ := http.NewRequest("GET", "/api/Machine/Get?id=nonexistent", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	// Should return not found
	if w4.Code != http.StatusNotFound && w4.Code != http.StatusInternalServerError {
		t.Errorf("Expected status NotFound or InternalServerError for nonexistent machine, got %d", w4.Code)
	}
}

func TestMachineHandler_CheckDeviceExist(t *testing.T) {
	router, _ := setupMachineTestRouter()

	// Test with valid deviceId
	req, _ := http.NewRequest("GET", "/api/Machine/CheckDeviceExist?deviceId=device123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status OK, BadRequest, or InternalServerError, got %d", w.Code)
	}

	// Test without deviceId parameter
	req2, _ := http.NewRequest("GET", "/api/Machine/CheckDeviceExist", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should return bad request due to missing deviceId
	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest for missing deviceId, got %d", w2.Code)
	}

	// Test with empty deviceId
	req3, _ := http.NewRequest("GET", "/api/Machine/CheckDeviceExist?deviceId=", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should return bad request for empty deviceId
	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest for empty deviceId, got %d", w3.Code)
	}
}

func TestMachineHandler_GetProductList(t *testing.T) {
	router, _ := setupMachineTestRouter()

	// Test with valid machineId
	req, _ := http.NewRequest("GET", "/api/Machine/GetProductList?machineId=machine123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status OK, BadRequest, or InternalServerError, got %d", w.Code)
	}

	// Test without machineId parameter
	req2, _ := http.NewRequest("GET", "/api/Machine/GetProductList", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should return bad request due to missing machineId
	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest for missing machineId, got %d", w2.Code)
	}

	// Test with empty machineId
	req3, _ := http.NewRequest("GET", "/api/Machine/GetProductList?machineId=", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should return bad request for empty machineId
	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest for empty machineId, got %d", w3.Code)
	}

	// Test with nonexistent machineId
	req4, _ := http.NewRequest("GET", "/api/Machine/GetProductList?machineId=nonexistent", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	// Should return OK with empty array or internal server error
	if w4.Code != http.StatusOK && w4.Code != http.StatusInternalServerError {
		t.Errorf("Expected status OK or InternalServerError for nonexistent machine, got %d", w4.Code)
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

	// Test with invalid JSON body
	req2, _ := http.NewRequest("POST", "/api/Machine/GetPaging", nil)
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should return forbidden due to lack of authentication (checked first)
	if w2.Code != http.StatusForbidden && w2.Code != http.StatusUnauthorized && w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status Forbidden, Unauthorized, or BadRequest, got %d", w2.Code)
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

	// Test with different request types
	req2, _ := http.NewRequest("GET", "/api/Machine/GetList?extra=param", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should still return forbidden or unauthorized regardless of extra params
	if w2.Code != http.StatusForbidden && w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Forbidden or Unauthorized with extra params, got %d", w2.Code)
	}
}

func TestMachineHandler_OpenOrCloseBusiness(t *testing.T) {
	router, _ := setupMachineTestRouter()

	// Test without authentication but with valid parameters
	req, _ := http.NewRequest("GET", "/api/Machine/OpenOrClose?id=machine123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return forbidden or unauthorized due to lack of authentication
	if w.Code != http.StatusForbidden && w.Code != http.StatusUnauthorized && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status Forbidden, Unauthorized, or BadRequest, got %d", w.Code)
	}

	// Test without id parameter
	req2, _ := http.NewRequest("GET", "/api/Machine/OpenOrClose", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Should return forbidden (auth checked first) or bad request
	if w2.Code != http.StatusBadRequest && w2.Code != http.StatusForbidden && w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status BadRequest, Forbidden, or Unauthorized, got %d", w2.Code)
	}

	// Test with empty id parameter
	req3, _ := http.NewRequest("GET", "/api/Machine/OpenOrClose?id=", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	// Should return forbidden (auth checked first) or bad request
	if w3.Code != http.StatusBadRequest && w3.Code != http.StatusForbidden && w3.Code != http.StatusUnauthorized {
		t.Errorf("Expected status BadRequest, Forbidden, or Unauthorized, got %d", w3.Code)
	}

	// Test with multiple parameters
	req4, _ := http.NewRequest("GET", "/api/Machine/OpenOrClose?id=machine123&extra=param", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	// Should still check authentication first
	if w4.Code != http.StatusForbidden && w4.Code != http.StatusUnauthorized && w4.Code != http.StatusBadRequest {
		t.Errorf("Expected status Forbidden, Unauthorized, or BadRequest with extra params, got %d", w4.Code)
	}
}
