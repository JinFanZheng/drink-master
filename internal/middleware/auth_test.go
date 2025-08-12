package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func setupTestAuth() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(JWTAuth())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "protected resource"})
	})
	return router
}

func generateTestJWT() (string, error) {
	claims := &JWTClaims{
		MemberID:       "test_member_123",
		MachineOwnerID: "test_owner_456",
		Role:           "Member",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_jwt_secret_change_this_in_production"
	}
	return token.SignedString([]byte(secret))
}

func TestJWTAuth_NoAuthorizationHeader(t *testing.T) {
	router := setupTestAuth()

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWTAuth_InvalidAuthorizationFormat(t *testing.T) {
	router := setupTestAuth()

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	router := setupTestAuth()

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token_string")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	router := setupTestAuth()

	token, err := generateTestJWT()
	if err != nil {
		t.Fatalf("Failed to generate test JWT: %v", err)
	}

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	token, err := generateTestJWT()
	if err != nil {
		t.Fatalf("Failed to generate test JWT: %v", err)
	}

	claims, err := validateJWT(token)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if claims.MemberID != "test_member_123" {
		t.Errorf("Expected MemberID to be 'test_member_123', got '%s'", claims.MemberID)
	}

	if claims.MachineOwnerID != "test_owner_456" {
		t.Errorf("Expected MachineOwnerID to be 'test_owner_456', got '%s'", claims.MachineOwnerID)
	}

	if claims.Role != "Member" {
		t.Errorf("Expected Role to be 'Member', got '%s'", claims.Role)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := validateJWT("invalid_token_string")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestGetCurrentMemberID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 测试没有member_id的情况
	memberID, exists := GetCurrentMemberID(c)
	if exists {
		t.Error("Expected GetCurrentMemberID to return false when no member_id is set")
	}
	if memberID != "" {
		t.Error("Expected memberID to be empty when not exists")
	}

	// 设置member_id并测试
	c.Set("member_id", "test_member_789")
	memberID, exists = GetCurrentMemberID(c)
	if !exists {
		t.Error("Expected GetCurrentMemberID to return true when member_id is set")
	}
	if memberID != "test_member_789" {
		t.Errorf("Expected memberID to be 'test_member_789', got '%s'", memberID)
	}
}

func TestGetCurrentMachineOwnerID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 测试没有machine_owner_id的情况
	ownerID, exists := GetCurrentMachineOwnerID(c)
	if exists {
		t.Error("Expected GetCurrentMachineOwnerID to return false when no machine_owner_id is set")
	}
	if ownerID != "" {
		t.Error("Expected ownerID to be empty when not exists")
	}

	// 设置machine_owner_id并测试
	c.Set("machine_owner_id", "test_owner_101112")
	ownerID, exists = GetCurrentMachineOwnerID(c)
	if !exists {
		t.Error("Expected GetCurrentMachineOwnerID to return true when machine_owner_id is set")
	}
	if ownerID != "test_owner_101112" {
		t.Errorf("Expected ownerID to be 'test_owner_101112', got '%s'", ownerID)
	}
}

func TestGetCurrentRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 测试没有role的情况
	role, exists := GetCurrentRole(c)
	if exists {
		t.Error("Expected GetCurrentRole to return false when no role is set")
	}
	if role != "" {
		t.Error("Expected role to be empty when not exists")
	}

	// 设置role并测试
	c.Set("role", "Owner")
	role, exists = GetCurrentRole(c)
	if !exists {
		t.Error("Expected GetCurrentRole to return true when role is set")
	}
	if role != "Owner" {
		t.Errorf("Expected role to be 'Owner', got '%s'", role)
	}
}

func TestIsMachineOwner(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 测试没有role的情况
	if IsMachineOwner(c) {
		t.Error("Expected IsMachineOwner to return false when no role is set")
	}

	// 设置非Owner角色
	c.Set("role", "Member")
	if IsMachineOwner(c) {
		t.Error("Expected IsMachineOwner to return false for Member role")
	}

	// 设置Owner角色
	c.Set("role", "Owner")
	if !IsMachineOwner(c) {
		t.Error("Expected IsMachineOwner to return true for Owner role")
	}
}
