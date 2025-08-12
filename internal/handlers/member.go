package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/repositories"
	"github.com/ddteam/drink-master/internal/services"
)

// MemberHandler 会员处理器 (对应MobileAPI MemberController)
type MemberHandler struct {
	*BaseHandler
	memberService *services.MemberService
}

// NewMemberHandler 创建会员处理器
func NewMemberHandler(db *gorm.DB) *MemberHandler {
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)
	memberService := services.NewMemberService(memberRepo, franchiseRepo)

	return &MemberHandler{
		BaseHandler:   NewBaseHandler(db),
		memberService: memberService,
	}
}

// Update 更新会员信息
// POST /api/Member/Update
func (h *MemberHandler) Update(c *gin.Context) {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	var req contracts.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 设置会员ID（从JWT获取）
	req.ID = memberID

	response, err := h.memberService.UpdateMember(memberID, req)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	h.SuccessResponse(c, response)
}

// AddFranchiseIntention 添加加盟意向
// POST /api/Member/AddFranchiseIntention
func (h *MemberHandler) AddFranchiseIntention(c *gin.Context) {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	var req contracts.CreateFranchiseIntentionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 设置会员ID（从JWT获取）
	req.MemberID = memberID

	response, err := h.memberService.CreateFranchiseIntention(memberID, req)
	if err != nil {
		// 检查是否是业务逻辑错误
		if err.Error() == "member already has pending franchise intention" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "已存在待处理的加盟意向",
				"code":    contracts.ErrorCodeFranchiseExists,
				"success": false,
			})
			return
		}

		h.InternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "加盟意向提交成功",
		"success": true,
	})
}

// GetUserInfo 获取用户信息（包含加盟意向）
// GET /api/Member/GetUserInfo
func (h *MemberHandler) GetUserInfo(c *gin.Context) {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	response, err := h.memberService.GetMemberInfo(memberID)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	h.SuccessResponse(c, response)
}
