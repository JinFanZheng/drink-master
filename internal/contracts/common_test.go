package contracts

import (
	"testing"
	"time"
)

func TestAPIResponse(t *testing.T) {
	response := APIResponse{
		Success: true,
		Data:    map[string]string{"test": "data"},
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}

	if response.Data == nil {
		t.Error("Expected Data to not be nil")
	}
}

func TestAPIError(t *testing.T) {
	now := time.Now()
	error := APIError{
		Code:      ErrorCodeValidation,
		Message:   "Test error",
		Timestamp: now,
		Path:      "/api/test",
		Method:    "POST",
		RequestID: "test-123",
	}

	if error.Code != ErrorCodeValidation {
		t.Errorf("Expected Code to be '%s', got '%s'", ErrorCodeValidation, error.Code)
	}
	if error.Message != "Test error" {
		t.Errorf("Expected Message to be 'Test error', got '%s'", error.Message)
	}
	if error.Timestamp != now {
		t.Error("Expected Timestamp to match")
	}
}

func TestErrorCodes(t *testing.T) {
	expectedCodes := []string{
		ErrorCodeValidation,
		ErrorCodeUnauthorized,
		ErrorCodeForbidden,
		ErrorCodeNotFound,
		ErrorCodeConflict,
		ErrorCodeInternalServer,
		ErrorCodeDatabaseError,
		ErrorCodeInvalidToken,
		ErrorCodeTokenExpired,
	}

	// 验证错误码不为空
	for _, code := range expectedCodes {
		if code == "" {
			t.Error("Error code should not be empty")
		}
	}

	// 验证错误码唯一性
	codeMap := make(map[string]bool)
	for _, code := range expectedCodes {
		if codeMap[code] {
			t.Errorf("Duplicate error code found: %s", code)
		}
		codeMap[code] = true
	}
}