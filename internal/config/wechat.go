package config

import "os"

// WeChatConfig represents WeChat configuration
type WeChatConfig struct {
	AppID     string
	AppSecret string
}

// NewWeChatConfig creates WeChat configuration from environment variables
func NewWeChatConfig() *WeChatConfig {
	return &WeChatConfig{
		AppID:     os.Getenv("WECHAT_APP_ID"),
		AppSecret: os.Getenv("WECHAT_APP_SECRET"),
	}
}
