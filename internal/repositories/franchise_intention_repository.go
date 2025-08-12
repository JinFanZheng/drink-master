package repositories

import (
	"fmt"

	"github.com/ddteam/drink-master/internal/models"
	"gorm.io/gorm"
)

// FranchiseIntentionRepository 加盟意向数据访问层
type FranchiseIntentionRepository struct {
	db *gorm.DB
}

// NewFranchiseIntentionRepository 创建加盟意向Repository实例
func NewFranchiseIntentionRepository(db *gorm.DB) *FranchiseIntentionRepository {
	return &FranchiseIntentionRepository{db: db}
}

// Create 创建加盟意向
func (r *FranchiseIntentionRepository) Create(intention *models.FranchiseIntention) error {
	err := r.db.Create(intention).Error
	if err != nil {
		return fmt.Errorf("failed to create franchise intention: %w", err)
	}
	return nil
}

// GetByID 根据ID获取加盟意向
func (r *FranchiseIntentionRepository) GetByID(id string) (*models.FranchiseIntention, error) {
	var intention models.FranchiseIntention
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Member").
		First(&intention).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("franchise intention not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get franchise intention: %w", err)
	}
	return &intention, nil
}

// GetByMemberID 根据会员ID获取加盟意向列表
func (r *FranchiseIntentionRepository) GetByMemberID(memberID string) ([]models.FranchiseIntention, error) {
	var intentions []models.FranchiseIntention
	err := r.db.Where("member_id = ? AND deleted_at IS NULL", memberID).
		Order("created_at DESC").
		Find(&intentions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get franchise intentions for member %s: %w", memberID, err)
	}
	return intentions, nil
}

// CheckExistingByMember 检查会员是否已有待处理的加盟意向
func (r *FranchiseIntentionRepository) CheckExistingByMember(memberID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.FranchiseIntention{}).
		Where("member_id = ? AND status = ? AND deleted_at IS NULL", memberID, "Pending").
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check existing franchise intention: %w", err)
	}
	return count > 0, nil
}

// Update 更新加盟意向
func (r *FranchiseIntentionRepository) Update(intention *models.FranchiseIntention) error {
	err := r.db.Save(intention).Error
	if err != nil {
		return fmt.Errorf("failed to update franchise intention: %w", err)
	}
	return nil
}

// UpdateStatus 更新加盟意向状态
func (r *FranchiseIntentionRepository) UpdateStatus(id string, status string) error {
	err := r.db.Model(&models.FranchiseIntention{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("status", status).Error
	if err != nil {
		return fmt.Errorf("failed to update franchise intention status: %w", err)
	}
	return nil
}

// Delete 软删除加盟意向
func (r *FranchiseIntentionRepository) Delete(id string) error {
	err := r.db.Where("id = ?", id).Delete(&models.FranchiseIntention{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete franchise intention: %w", err)
	}
	return nil
}

// GetPaginated 分页获取加盟意向列表
func (r *FranchiseIntentionRepository) GetPaginated(
	offset, limit int, status string,
) ([]models.FranchiseIntention, int64, error) {
	var intentions []models.FranchiseIntention
	var total int64

	query := r.db.Model(&models.FranchiseIntention{}).Where("deleted_at IS NULL")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count franchise intentions: %w", err)
	}

	// 获取分页数据
	err = query.Preload("Member").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&intentions).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated franchise intentions: %w", err)
	}

	return intentions, total, nil
}
