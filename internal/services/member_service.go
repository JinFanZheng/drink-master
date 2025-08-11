package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/models"
	"github.com/ddteam/drink-master/internal/repositories"
)

// MemberService 会员业务逻辑层
type MemberService struct {
	memberRepo             *repositories.MemberRepository
	franchiseIntentionRepo *repositories.FranchiseIntentionRepository
}

// NewMemberService 创建会员Service实例
func NewMemberService(
	memberRepo *repositories.MemberRepository,
	franchiseIntentionRepo *repositories.FranchiseIntentionRepository,
) *MemberService {
	return &MemberService{
		memberRepo:             memberRepo,
		franchiseIntentionRepo: franchiseIntentionRepo,
	}
}

// UpdateMember 更新会员信息
func (s *MemberService) UpdateMember(
	memberID string,
	req contracts.UpdateMemberRequest,
) (*contracts.UpdateMemberResponse, error) {
	// 获取现有会员信息
	member, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	// 验证输入数据
	if req.Nickname == "" {
		return nil, fmt.Errorf("nickname cannot be empty")
	}
	if req.Avatar == "" {
		return nil, fmt.Errorf("avatar cannot be empty")
	}

	// 更新会员信息
	member.Nickname = req.Nickname
	member.Avatar = req.Avatar
	member.UpdatedAt = time.Now()

	err = s.memberRepo.Update(member)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	// 返回更新结果
	response := &contracts.UpdateMemberResponse{
		ID:       member.ID,
		Nickname: member.Nickname,
		Avatar:   member.Avatar,
		Success:  true,
	}

	return response, nil
}

// CreateFranchiseIntention 创建加盟意向
func (s *MemberService) CreateFranchiseIntention(
	memberID string,
	req contracts.CreateFranchiseIntentionRequest,
) (*contracts.CreateFranchiseIntentionResponse, error) {
	// 验证会员是否存在
	_, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %w", err)
	}

	// 检查是否已有待处理的加盟意向
	hasExisting, err := s.franchiseIntentionRepo.CheckExistingByMember(memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing franchise intention: %w", err)
	}
	if hasExisting {
		return nil, fmt.Errorf("member already has pending franchise intention")
	}

	// 验证输入数据
	if req.ContactName == "" {
		return nil, fmt.Errorf("contact name cannot be empty")
	}
	if req.ContactPhone == "" {
		return nil, fmt.Errorf("contact phone cannot be empty")
	}
	if req.IntendedLocation == "" {
		return nil, fmt.Errorf("intended location cannot be empty")
	}

	// 创建加盟意向
	intention := &models.FranchiseIntention{
		ID:               s.generateFranchiseIntentionID(),
		MemberID:         memberID,
		ContactName:      req.ContactName,
		ContactPhone:     req.ContactPhone,
		IntendedLocation: req.IntendedLocation,
		Remarks:          req.Remarks,
		Status:           "Pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.franchiseIntentionRepo.Create(intention)
	if err != nil {
		return nil, fmt.Errorf("failed to create franchise intention: %w", err)
	}

	// 返回创建结果
	response := &contracts.CreateFranchiseIntentionResponse{
		ID:               intention.ID,
		MemberID:         intention.MemberID,
		ContactName:      intention.ContactName,
		ContactPhone:     intention.ContactPhone,
		IntendedLocation: intention.IntendedLocation,
		Status:           intention.Status,
		CreatedAt:        intention.CreatedAt,
		Success:          true,
	}

	return response, nil
}

// GetMemberInfo 获取会员详细信息
func (s *MemberService) GetMemberInfo(memberID string) (*contracts.GetMemberInfoResponse, error) {
	// 获取会员信息和加盟意向
	member, intentions, err := s.memberRepo.GetMemberWithFranchiseIntentions(memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member info: %w", err)
	}

	// 转换加盟意向为摘要格式
	var intentionSummaries []contracts.FranchiseIntentionSummary
	for _, intention := range intentions {
		summary := contracts.FranchiseIntentionSummary{
			ID:               intention.ID,
			ContactName:      intention.ContactName,
			IntendedLocation: intention.IntendedLocation,
			Status:           contracts.FranchiseIntentionStatus(intention.Status),
			CreatedAt:        intention.CreatedAt,
		}
		intentionSummaries = append(intentionSummaries, summary)
	}

	// 构建响应
	response := &contracts.GetMemberInfoResponse{
		ID:                  member.ID,
		Nickname:            member.Nickname,
		Avatar:              member.Avatar,
		WeChatOpenID:        member.WeChatOpenId,
		Role:                member.Role,
		IsAdmin:             member.IsAdmin,
		CreatedAt:           member.CreatedAt,
		UpdatedAt:           member.UpdatedAt,
		FranchiseIntentions: intentionSummaries,
	}

	// 设置机主ID（如果存在）
	if member.MachineOwnerId != nil {
		response.MachineOwnerID = *member.MachineOwnerId
	}

	return response, nil
}

// ValidateMemberExists 验证会员是否存在
func (s *MemberService) ValidateMemberExists(memberID string) error {
	_, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		return fmt.Errorf("member validation failed: %w", err)
	}
	return nil
}

// UpdateFranchiseIntentionStatus 更新加盟意向状态（管理员功能）
func (s *MemberService) UpdateFranchiseIntentionStatus(intentionID string, status string) error {
	// 验证状态值
	validStatuses := []string{"Pending", "Approved", "Rejected"}
	statusValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			statusValid = true
			break
		}
	}
	if !statusValid {
		return fmt.Errorf("invalid status: %s", status)
	}

	// 更新状态
	err := s.franchiseIntentionRepo.UpdateStatus(intentionID, status)
	if err != nil {
		return fmt.Errorf("failed to update franchise intention status: %w", err)
	}

	return nil
}

// generateFranchiseIntentionID 生成加盟意向ID
func (s *MemberService) generateFranchiseIntentionID() string {
	// 生成8位随机字符串
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return "fi-" + string(b)
}
