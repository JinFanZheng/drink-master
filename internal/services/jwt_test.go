package services

import (
	"os"
	"testing"

	"github.com/ddteam/drink-master/internal/models"
)

func TestNewJWTService(t *testing.T) {
	service := NewJWTService()

	if service == nil {
		t.Error("JWT service should not be nil")
	}

	if len(service.secret) == 0 {
		t.Error("JWT secret should not be empty")
	}
}

func TestGenerateToken(t *testing.T) {
	service := NewJWTService()

	member := &models.Member{
		ID:       "test_member_id",
		Nickname: "Test User",
		Role:     "Member",
	}

	token, err := service.GenerateToken(member)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	service := NewJWTService()

	member := &models.Member{
		ID:       "test_member_id",
		Nickname: "Test User",
		Role:     "Member",
	}

	// Generate token
	token, err := service.GenerateToken(member)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if claims.MemberID != member.ID {
		t.Errorf("Expected member ID %s, got %s", member.ID, claims.MemberID)
	}

	if claims.Role != member.Role {
		t.Errorf("Expected role %s, got %s", member.Role, claims.Role)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	service := NewJWTService()

	_, err := service.ValidateToken("invalid_token")

	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestJWTService_WithCustomSecret(t *testing.T) {
	// Set custom secret
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "custom_test_secret")
	defer os.Setenv("JWT_SECRET", originalSecret)

	service := NewJWTService()

	if string(service.secret) != "custom_test_secret" {
		t.Errorf("Expected custom secret, got %s", string(service.secret))
	}
}
