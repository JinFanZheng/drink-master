package wechat

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	appID := "test_app_id"
	appSecret := "test_app_secret"

	client := NewClient(appID, appSecret)

	if client.appID != appID {
		t.Errorf("Expected appID %s, got %s", appID, client.appID)
	}

	if client.appSecret != appSecret {
		t.Errorf("Expected appSecret %s, got %s", appSecret, client.appSecret)
	}

	if client.client == nil {
		t.Error("HTTP client should not be nil")
	}
}

func TestJsCode2Session_EmptyCode(t *testing.T) {
	client := NewClient("test_app_id", "test_app_secret")

	_, err := client.JsCode2Session("")

	if err == nil {
		t.Error("Expected error for empty code")
	}

	if err.Error() != "code cannot be empty" {
		t.Errorf("Expected 'code cannot be empty' error, got %s", err.Error())
	}
}
