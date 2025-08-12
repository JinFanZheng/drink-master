package contracts

import "time"

// Member contracts 会员管理相关的API契约

// UpdateMemberRequest 更新会员信息请求
type UpdateMemberRequest struct {
	ID       string `json:"id,omitempty"` // 会员ID (从JWT获取，请求中不需要)
	Nickname string `json:"nickname" binding:"required" validate:"min=2,max=50" example:"张三"`
	Avatar   string `json:"avatar" binding:"required" validate:"url" example:"https://example.com/avatar.jpg"`
}

// UpdateMemberResponse 更新会员信息响应
type UpdateMemberResponse struct {
	ID       string `json:"id" example:"member-123"`
	Nickname string `json:"nickname" example:"张三"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Success  bool   `json:"success" example:"true"`
}

// CreateFranchiseIntentionRequest 创建加盟意向请求
type CreateFranchiseIntentionRequest struct {
	MemberID         string `json:"memberId,omitempty"` // 会员ID (从JWT获取)
	ContactName      string `json:"contactName" binding:"required" validate:"min=2,max=50" example:"张三"`
	ContactPhone     string `json:"contactPhone" binding:"required" validate:"phone" example:"13800138000"`
	IntendedLocation string `json:"intendedLocation" binding:"required" validate:"min=5,max=200" example:"北京市朝阳区建国门外大街"`
	Remarks          string `json:"remarks,omitempty" validate:"max=500" example:"希望能在繁华地段开店"`
}

// CreateFranchiseIntentionResponse 创建加盟意向响应
type CreateFranchiseIntentionResponse struct {
	ID               string    `json:"id" example:"intention-123"`
	MemberID         string    `json:"memberId" example:"member-123"`
	ContactName      string    `json:"contactName" example:"张三"`
	ContactPhone     string    `json:"contactPhone" example:"13800138000"`
	IntendedLocation string    `json:"intendedLocation" example:"北京市朝阳区建国门外大街"`
	Status           string    `json:"status" example:"Pending"`
	CreatedAt        time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	Success          bool      `json:"success" example:"true"`
}

// FranchiseIntentionStatus 加盟意向状态枚举
type FranchiseIntentionStatus string

const (
	FranchiseIntentionStatusPending  FranchiseIntentionStatus = "Pending"  // 待处理
	FranchiseIntentionStatusApproved FranchiseIntentionStatus = "Approved" // 已批准
	FranchiseIntentionStatusRejected FranchiseIntentionStatus = "Rejected" // 已拒绝
)

// GetMemberInfoResponse 获取会员信息响应
type GetMemberInfoResponse struct {
	ID                  string                      `json:"id" example:"member-123"`
	Nickname            string                      `json:"nickname" example:"张三"`
	Avatar              string                      `json:"avatar" example:"https://example.com/avatar.jpg"`
	WeChatOpenID        string                      `json:"weChatOpenId" example:"wx_openid_123"`
	Role                string                      `json:"role" example:"Member"`
	MachineOwnerID      string                      `json:"machineOwnerId,omitempty" example:"owner-123"`
	IsAdmin             bool                        `json:"isAdmin" example:"false"`
	CreatedAt           time.Time                   `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt           time.Time                   `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
	FranchiseIntentions []FranchiseIntentionSummary `json:"franchiseIntentions,omitempty"`
}

// FranchiseIntentionSummary 加盟意向摘要
type FranchiseIntentionSummary struct {
	ID               string                   `json:"id" example:"intention-123"`
	ContactName      string                   `json:"contactName" example:"张三"`
	IntendedLocation string                   `json:"intendedLocation" example:"北京市朝阳区"`
	Status           FranchiseIntentionStatus `json:"status" example:"Pending"`
	CreatedAt        time.Time                `json:"createdAt" example:"2023-01-01T00:00:00Z"`
}

// Member相关的错误码
const (
	ErrorCodeMemberNotFound      = "MEMBER_NOT_FOUND"
	ErrorCodeMemberUpdateFailed  = "MEMBER_UPDATE_FAILED"
	ErrorCodeFranchiseExists     = "FRANCHISE_INTENTION_EXISTS"
	ErrorCodeFranchiseCreateFail = "FRANCHISE_CREATE_FAILED"
	ErrorCodeInvalidMemberData   = "INVALID_MEMBER_DATA"
)
