package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MemberHandler 会员处理器 (对应MobileAPI MemberController)
type MemberHandler struct {
	*BaseHandler
}

// NewMemberHandler 创建会员处理器
func NewMemberHandler(db *gorm.DB) *MemberHandler {
	return &MemberHandler{
		BaseHandler: NewBaseHandler(db),
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

	// TODO: 实现会员信息更新逻辑
	h.SuccessResponseWithMessage(c, map[string]interface{}{
		"id": memberID,
	}, "会员信息更新成功")
}

// AddFranchiseIntention 添加加盟意向
// POST /api/Member/AddFranchiseIntention
func (h *MemberHandler) AddFranchiseIntention(c *gin.Context) {
	memberID, exists := h.GetMemberID(c)
	if !exists {
		h.UnauthorizedResponse(c, "无效的用户信息")
		return
	}

	// TODO: 实现加盟意向添加逻辑
	h.SuccessResponseWithMessage(c, map[string]interface{}{
		"memberId": memberID,
	}, "加盟意向提交成功")
}
