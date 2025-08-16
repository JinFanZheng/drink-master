package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

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
	member.Nickname = &req.Nickname
	member.Avatar = &req.Avatar
	// UpdatedAt 将在数据库中自动更新，不需要手动设置

	err = s.memberRepo.Update(member)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	// 返回更新结果
	nickname := ""
	avatar := ""
	if member.Nickname != nil {
		nickname = *member.Nickname
	}
	if member.Avatar != nil {
		avatar = *member.Avatar
	}

	response := &contracts.UpdateMemberResponse{
		ID:       member.ID,
		Nickname: nickname,
		Avatar:   avatar,
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
	company := req.ContactName // Use ContactName as Company
	name := req.ContactName
	mobile := req.ContactPhone
	area := req.IntendedLocation
	remark := req.Remarks

	intention := &models.FranchiseIntention{
		ID:        s.generateFranchiseIntentionID(),
		MemberId:  &memberID,
		Company:   &company,
		Name:      &name,
		Mobile:    &mobile,
		Area:      &area,
		Remark:    &remark,
		IsHandled: models.BitBool(0),
		Version:   0,
		CreatedOn: time.Now(),
	}

	err = s.franchiseIntentionRepo.Create(intention)
	if err != nil {
		return nil, fmt.Errorf("failed to create franchise intention: %w", err)
	}

	// 返回创建结果
	memberId := ""
	contactName := ""
	contactPhone := ""
	intendedLocation := ""

	if intention.MemberId != nil {
		memberId = *intention.MemberId
	}
	if intention.Name != nil {
		contactName = *intention.Name
	}
	if intention.Mobile != nil {
		contactPhone = *intention.Mobile
	}
	if intention.Area != nil {
		intendedLocation = *intention.Area
	}

	// Convert IsHandled status to proper status string
	status := "Pending"
	if intention.IsHandled.Bool() {
		status = "Handled"
	}

	response := &contracts.CreateFranchiseIntentionResponse{
		ID:               intention.ID,
		MemberID:         memberId,
		ContactName:      contactName,
		ContactPhone:     contactPhone,
		IntendedLocation: intendedLocation,
		Status:           status, // Use proper status string instead of bool conversion
		CreatedAt:        intention.CreatedOn,
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
		contactName := ""
		intendedLocation := ""

		if intention.Name != nil {
			contactName = *intention.Name
		}
		if intention.Area != nil {
			intendedLocation = *intention.Area
		}

		// Convert IsHandled status to proper status string
		intentionStatus := "Pending"
		if intention.IsHandled.Bool() {
			intentionStatus = "Handled"
		}

		summary := contracts.FranchiseIntentionSummary{
			ID:               intention.ID,
			ContactName:      contactName,
			IntendedLocation: intendedLocation,
			Status:           contracts.FranchiseIntentionStatus(intentionStatus),
			CreatedAt:        intention.CreatedOn,
		}
		intentionSummaries = append(intentionSummaries, summary)
	}

	// 构建响应
	nickname := ""
	avatar := ""
	wechatOpenId := ""

	if member.Nickname != nil {
		nickname = *member.Nickname
	}
	if member.Avatar != nil {
		avatar = *member.Avatar
	}
	if member.WeChatOpenId != nil {
		wechatOpenId = *member.WeChatOpenId
	}

	response := &contracts.GetMemberInfoResponse{
		ID:                  member.ID,
		Nickname:            nickname,
		Avatar:              avatar,
		WeChatOpenID:        wechatOpenId,
		Role:                fmt.Sprintf("%d", member.Role), // Convert int to string
		IsAdmin:             member.IsAdmin.Bool(),
		CreatedAt:           member.CreatedOn,
		UpdatedAt:           member.CreatedOn, // Use CreatedOn since UpdatedOn might be nil
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

	// 更新状态 - 将字符串状态转换为bool
	isHandled := status != "Pending" // Pending表示未处理，其他状态都表示已处理
	err := s.franchiseIntentionRepo.UpdateStatus(intentionID, isHandled)
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

// 兼容旧版本的构造函数（用于PR #18）
func NewMemberServiceCompat(db *gorm.DB) *MemberService {
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)
	return NewMemberService(memberRepo, franchiseRepo)
}

// FindByOpenID finds a member by WeChat OpenID (for PR #18 compatibility)
func (s *MemberService) FindByOpenID(openID string) (*models.Member, error) {
	return s.memberRepo.GetByWeChatOpenID(openID)
}

// FindOrCreateByOpenID finds a member by OpenID or creates a new one (for PR #18 compatibility)
func (s *MemberService) FindOrCreateByOpenID(openID, nickname, avatarUrl string) (*models.Member, error) {
	// Try to find existing member
	member, err := s.memberRepo.GetByWeChatOpenID(openID)
	if err == nil {
		// Member exists, update avatar and nickname if provided
		if nickname != "" {
			member.Nickname = &nickname
		}
		if avatarUrl != "" {
			member.Avatar = &avatarUrl
		}
		if updateErr := s.memberRepo.Update(member); updateErr != nil {
			return nil, updateErr
		}
		return member, nil
	}

	// Member not found, create new one (check for "member not found" in error message)
	if err != nil && (err == gorm.ErrRecordNotFound ||
		fmt.Sprintf("%v", err) == fmt.Sprintf("member not found with openID: %s", openID)) {
		newMember := &models.Member{
			ID:           uuid.New().String(),
			Nickname:     &nickname,
			Avatar:       &avatarUrl,
			WeChatOpenId: &openID,
			Role:         1, // Default role as int (1 for Member)
			Version:      0,
			CreatedOn:    time.Now(),
		}

		if createErr := s.memberRepo.Create(newMember); createErr != nil {
			return nil, createErr
		}
		return newMember, nil
	}

	return nil, err
}

// FindByID finds a member by ID (for PR #18 compatibility)
func (s *MemberService) FindByID(id string) (*models.Member, error) {
	return s.memberRepo.GetByID(id)
}
