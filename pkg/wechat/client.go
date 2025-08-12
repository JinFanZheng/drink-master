package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client represents WeChat API client
type Client struct {
	appID     string
	appSecret string
	client    *http.Client
}

// NewClient creates a new WeChat client
func NewClient(appID, appSecret string) *Client {
	return &Client{
		appID:     appID,
		appSecret: appSecret,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SessionResponse represents WeChat jscode2session response
type SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}

// JsCode2Session exchanges code for user session information
func (c *Client) JsCode2Session(code string) (*SessionResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("code cannot be empty")
	}

	// Build request URL
	params := url.Values{}
	params.Add("appid", c.appID)
	params.Add("secret", c.appSecret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")

	requestURL := "https://api.weixin.qq.com/sns/jscode2session?" + params.Encode()

	// Make HTTP request
	resp, err := c.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to request WeChat API: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var sessionResp SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for WeChat API errors
	if sessionResp.ErrCode != 0 {
		return nil, fmt.Errorf("WeChat API error: %d - %s", sessionResp.ErrCode, sessionResp.ErrMsg)
	}

	return &sessionResp, nil
}
