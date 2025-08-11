package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/middleware"
)

// BaseHandler 基础控制器结构 (对应MobileAPI BaseController)
type BaseHandler struct {
	db *gorm.DB
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(db *gorm.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

// GetMemberID 从JWT token获取用户ID
func (h *BaseHandler) GetMemberID(c *gin.Context) (string, bool) {
	return middleware.GetCurrentMemberID(c)
}

// GetMachineOwnerID 从JWT token获取机主ID
func (h *BaseHandler) GetMachineOwnerID(c *gin.Context) (string, bool) {
	return middleware.GetCurrentMachineOwnerID(c)
}

// IsMachineOwner 检查用户角色是否为Owner
func (h *BaseHandler) IsMachineOwner(c *gin.Context) bool {
	return middleware.IsMachineOwner(c)
}

// GetCurrentRole 获取当前用户角色
func (h *BaseHandler) GetCurrentRole(c *gin.Context) (string, bool) {
	return middleware.GetCurrentRole(c)
}

// SuccessResponse 返回成功响应
func (h *BaseHandler) SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, contracts.APIResponse{
		Success: true,
		Data:    data,
		Meta: &contracts.Meta{
			Timestamp: time.Now(),
			RequestID: getRequestID(c),
		},
	})
}

// SuccessResponseWithMessage 返回带消息的成功响应
func (h *BaseHandler) SuccessResponseWithMessage(c *gin.Context, data interface{}, message string) {
	response := map[string]interface{}{
		"success": true,
		"data":    data,
		"message": message,
	}
	c.JSON(http.StatusOK, response)
}

// ErrorResponse 返回错误响应
func (h *BaseHandler) ErrorResponse(c *gin.Context, code int, errCode string, message string) {
	c.JSON(code, contracts.APIResponse{
		Success: false,
		Error: &contracts.APIError{
			Code:      errCode,
			Message:   message,
			Timestamp: time.Now(),
			Path:      c.Request.URL.Path,
			Method:    c.Request.Method,
			RequestID: getRequestID(c),
		},
	})
}

// ValidationErrorResponse 返回验证错误响应
func (h *BaseHandler) ValidationErrorResponse(c *gin.Context, err error) {
	h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, err.Error())
}

// NotFoundResponse 返回404响应
func (h *BaseHandler) NotFoundResponse(c *gin.Context, message string) {
	h.ErrorResponse(c, http.StatusNotFound, contracts.ErrorCodeNotFound, message)
}

// UnauthorizedResponse 返回401响应
func (h *BaseHandler) UnauthorizedResponse(c *gin.Context, message string) {
	h.ErrorResponse(c, http.StatusUnauthorized, contracts.ErrorCodeUnauthorized, message)
}

// ForbiddenResponse 返回403响应
func (h *BaseHandler) ForbiddenResponse(c *gin.Context, message string) {
	h.ErrorResponse(c, http.StatusForbidden, contracts.ErrorCodeForbidden, message)
}

// InternalErrorResponse 返回500响应
func (h *BaseHandler) InternalErrorResponse(c *gin.Context, err error) {
	h.ErrorResponse(c, http.StatusInternalServerError, contracts.ErrorCodeInternalServer, "内部服务器错误")
}

// PagingResponse 返回分页响应 (基于MobileAPI PagingResponse)
func (h *BaseHandler) PagingResponse(c *gin.Context, items interface{}, totalCount int64, pageIndex, pageSize int) {
	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"items":      items,
			"totalCount": totalCount,
			"pageIndex":  pageIndex,
			"pageSize":   pageSize,
		},
	}
	c.JSON(http.StatusOK, response)
}

// getRequestID 获取请求ID的辅助函数
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}
