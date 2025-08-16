package repositories

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

// MemberRepository 会员数据访问层
type MemberRepository struct {
	db *gorm.DB
}

// NewMemberRepository 创建会员Repository实例
func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

// GetByID 根据ID获取会员信息
func (r *MemberRepository) GetByID(id string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("id = ?", id).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("member not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get member: %w", err)
	}
	return &member, nil
}

// GetByWeChatOpenID 根据微信OpenID获取会员信息
func (r *MemberRepository) GetByWeChatOpenID(openID string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("WeChatOpenId = ?", openID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("member not found with openID: %s", openID)
		}
		return nil, fmt.Errorf("failed to get member by openID: %w", err)
	}
	return &member, nil
}

// Update 更新会员信息
func (r *MemberRepository) Update(member *models.Member) error {
	err := r.db.Save(member).Error
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}
	return nil
}

// Create 创建新会员
func (r *MemberRepository) Create(member *models.Member) error {
	err := r.db.Create(member).Error
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}
	return nil
}

// Delete 软删除会员
func (r *MemberRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.Member{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}
	return nil
}

// GetMemberWithFranchiseIntentions 获取会员及其加盟意向
func (r *MemberRepository) GetMemberWithFranchiseIntentions(
	memberID string,
) (*models.Member, []models.FranchiseIntention, error) {
	var member models.Member
	err := r.db.Where("id = ?", memberID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("member not found: %s", memberID)
		}
		return nil, nil, fmt.Errorf("failed to get member: %w", err)
	}

	var intentions []models.FranchiseIntention
	err = r.db.Where("MemberId = ?", memberID).
		Order("CreatedOn DESC").
		Find(&intentions).Error
	if err != nil {
		return &member, nil, fmt.Errorf("failed to get franchise intentions: %w", err)
	}

	return &member, intentions, nil
}
