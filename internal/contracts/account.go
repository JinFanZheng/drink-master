package contracts

// CheckUserInfoRequest 检查用户信息请求
type CheckUserInfoRequest struct {
	Code  string `form:"code" binding:"required" example:"wx_js_code"`
	AppId string `form:"appId" example:"wx_app_id"`
}

// CheckUserInfoResponse 检查用户信息响应
type CheckUserInfoResponse struct {
	Id             string `json:"id" example:"member_123"`
	AvatarUrl      string `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	Nickname       string `json:"nickname" example:"用户昵称"`
	IsMachineOwner bool   `json:"isMachineOwner" example:"false"`
	Token          string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// WeChatLoginRequest 微信登录请求
type WeChatLoginRequest struct {
	AppId     string `json:"appId" example:"wx_app_id"`
	Code      string `json:"code" binding:"required" example:"wx_js_code"`
	AvatarUrl string `json:"avatarUrl" binding:"required" example:"https://example.com/avatar.jpg"`
	NickName  string `json:"nickName" binding:"required" example:"用户昵称"`
}

// WeChatLoginResponse 微信登录响应
type WeChatLoginResponse struct {
	Id             string `json:"id" example:"member_123"`
	AvatarUrl      string `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	Nickname       string `json:"nickname" example:"用户昵称"`
	IsMachineOwner bool   `json:"isMachineOwner" example:"false"`
	Token          string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// GetUserInfoResponse 获取用户信息响应
type GetUserInfoResponse struct {
	Id             string `json:"id" example:"member_123"`
	AvatarUrl      string `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	Nickname       string `json:"nickname" example:"用户昵称"`
	IsMachineOwner bool   `json:"isMachineOwner" example:"false"`
}
