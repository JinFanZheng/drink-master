package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AccountHandler 账户处理器 (对应MobileAPI AccountController)
type AccountHandler struct {
	*BaseHandler
}

// NewAccountHandler 创建账户处理器
func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

// CheckUserInfo 检查用户信息
// GET /api/Account/CheckUserInfo?code=wx_code&appId=wx_app_id
func (h *AccountHandler) CheckUserInfo(c *gin.Context) {
	// 验证必需的openId参数
	openId := c.Query("openId")
	if openId == "" {
		h.ValidationErrorResponse(c, fmt.Errorf("openId parameter is required"))
		return
	}

	// TODO: 实现微信code验证逻辑
	// 临时返回成功响应
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":             "temp_user_id",
			"avatarUrl":      "",
			"nickname":       "临时用户",
			"isMachineOwner": false,
			"token":          "temp_token",
		},
	}
	c.JSON(http.StatusOK, response)
}

// WeChatLogin 微信登录
// POST /api/Account/WeChatLogin
func (h *AccountHandler) WeChatLogin(c *gin.Context) {
	// 验证JSON格式
	var loginRequest map[string]interface{}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		h.ValidationErrorResponse(c, fmt.Errorf("invalid JSON format: %v", err))
		return
	}

	// TODO: 实现微信登录逻辑
	// 临时返回成功响应
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":             "temp_user_id",
			"avatarUrl":      "",
			"nickname":       "临时用户",
			"isMachineOwner": false,
			"token":          "temp_token",
		},
	}
	c.JSON(http.StatusOK, response)
}

// CheckLogin 检查登录状态
// GET /api/Account/CheckLogin (需要Authorization Bearer token)
func (h *AccountHandler) CheckLogin(c *gin.Context) {
	// 如果能到达这里，说明JWT验证已通过
	c.String(http.StatusOK, "ok")
}

// GetUserInfo 获取用户信息
// GET /api/Account/GetUserInfo (需要Authorization Bearer token)
func (h *AccountHandler) GetUserInfo(c *gin.Context) {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	// TODO: 从数据库获取用户信息
	// 临时返回成功响应
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":             memberID,
			"avatarUrl":      "",
			"nickname":       "用户",
			"isMachineOwner": h.IsMachineOwner(c),
		},
	}
	c.JSON(http.StatusOK, response)
}
