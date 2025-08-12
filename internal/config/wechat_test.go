package config

import (
	"os"
	"testing"
)

func TestNewWeChatConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("WECHAT_APP_ID", "test-app-id")
	os.Setenv("WECHAT_APP_SECRET", "test-app-secret")
	defer func() {
		os.Unsetenv("WECHAT_APP_ID")
		os.Unsetenv("WECHAT_APP_SECRET")
	}()

	config := NewWeChatConfig()

	if config == nil {
		t.Fatal("expected config to be created")
	}

	if config.AppID != "test-app-id" {
		t.Errorf("expected AppID 'test-app-id', got '%s'", config.AppID)
	}

	if config.AppSecret != "test-app-secret" {
		t.Errorf("expected AppSecret 'test-app-secret', got '%s'", config.AppSecret)
	}
}

func TestNewWeChatConfig_DefaultValues(t *testing.T) {
	// Ensure environment variables are not set
	os.Unsetenv("WECHAT_APP_ID")
	os.Unsetenv("WECHAT_APP_SECRET")

	config := NewWeChatConfig()

	if config == nil {
		t.Fatal("expected config to be created")
	}

	// Should use default/empty values when env vars are not set
	if config.AppID != "" {
		t.Errorf("expected empty AppID, got '%s'", config.AppID)
	}

	if config.AppSecret != "" {
		t.Errorf("expected empty AppSecret, got '%s'", config.AppSecret)
	}
}
