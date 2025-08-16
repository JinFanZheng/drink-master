package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/services"
	"github.com/ddteam/drink-master/pkg/wechat"
)

// Helper function to get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// AccountHandler 账户处理器 (对应MobileAPI AccountController)
type AccountHandler struct {
	*BaseHandler
	memberService *services.MemberService
	jwtService    *services.JWTService
	cacheManager  *services.CacheManager
	wechatClient  *wechat.Client
}

// NewAccountHandler 创建账户处理器
func NewAccountHandler(db *gorm.DB, wechatClient *wechat.Client) *AccountHandler {
	return &AccountHandler{
		BaseHandler:   NewBaseHandler(db),
		memberService: services.NewMemberServiceCompat(db),
		jwtService:    services.NewJWTService(),
		cacheManager:  services.NewCacheManager(),
		wechatClient:  wechatClient,
	}
}

// CheckUserInfo 检查用户信息
// @Summary 检查用户信息
// @Description 通过微信code检查用户信息
// @Tags Account
// @Accept json
// @Produce json
// @Param code query string true "微信授权码"
// @Param appId query string true "微信AppID"
// @Success 200 {object} contracts.CheckUserInfoResponse
// @Failure 400 {object} contracts.APIResponse
// @Router /Account/CheckUserInfo [get]
func (h *AccountHandler) CheckUserInfo(c *gin.Context) {
	var req contracts.CheckUserInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 验证微信code
	session, err := h.wechatClient.JsCode2Session(req.Code)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, "微信验证失败: "+err.Error())
		return
	}

	// 查找用户
	member, err := h.memberService.FindByOpenID(session.OpenID)
	if err != nil {
		// 用户不存在，返回空用户信息
		response := map[string]interface{}{
			"success": true,
			"data": contracts.CheckUserInfoResponse{
				Id:             "",
				AvatarUrl:      "",
				Nickname:       "",
				IsMachineOwner: false,
				Token:          "",
			},
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// 生成token
	token, err := h.jwtService.GenerateToken(member)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// 缓存登录状态
	h.cacheManager.SetLoginStatus(member.ID, token)

	// 返回用户信息
	response := map[string]interface{}{
		"success": true,
		"data": contracts.CheckUserInfoResponse{
			Id:             member.ID,
			AvatarUrl:      getStringValue(member.Avatar),
			Nickname:       getStringValue(member.Nickname),
			IsMachineOwner: member.Role == 2, // 2 for Owner role
			Token:          token,
		},
	}
	c.JSON(http.StatusOK, response)
}

// WeChatLogin 微信登录
// @Summary 微信登录
// @Description 通过微信授权码进行用户登录
// @Tags Account
// @Accept json
// @Produce json
// @Param request body contracts.WeChatLoginRequest true "登录请求"
// @Success 200 {object} contracts.WeChatLoginResponse
// @Failure 400 {object} contracts.APIResponse
// @Router /Account/WeChatLogin [post]
func (h *AccountHandler) WeChatLogin(c *gin.Context) {
	var req contracts.WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 验证微信code
	session, err := h.wechatClient.JsCode2Session(req.Code)
	if err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, "微信验证失败: "+err.Error())
		return
	}

	// 查找或创建用户
	member, err := h.memberService.FindOrCreateByOpenID(session.OpenID, req.NickName, req.AvatarUrl)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// 生成token
	token, err := h.jwtService.GenerateToken(member)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// 缓存登录状态
	h.cacheManager.SetLoginStatus(member.ID, token)

	// 返回用户信息
	response := map[string]interface{}{
		"success": true,
		"data": contracts.WeChatLoginResponse{
			Id:             member.ID,
			AvatarUrl:      getStringValue(member.Avatar),
			Nickname:       getStringValue(member.Nickname),
			IsMachineOwner: member.Role == 2, // 2 for Owner role
			Token:          token,
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

	// 从数据库获取用户信息
	member, err := h.memberService.FindByID(memberID)
	if err != nil {
		h.NotFoundResponse(c, "用户不存在")
		return
	}

	// 返回用户信息
	response := map[string]interface{}{
		"success": true,
		"data": contracts.GetUserInfoResponse{
			Id:             member.ID,
			AvatarUrl:      getStringValue(member.Avatar),
			Nickname:       getStringValue(member.Nickname),
			IsMachineOwner: member.Role == 2, // 2 for Owner role
		},
	}
	c.JSON(http.StatusOK, response)
}
