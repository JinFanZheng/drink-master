package repositories

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

// MachineRepository 售货机仓储接口
type MachineRepositoryInterface interface {
	GetByID(id string) (*models.Machine, error)
	GetByDeviceID(deviceID string) (*models.Machine, error)
	GetList(machineOwnerID string) ([]*models.Machine, error)
	GetPaging(machineOwnerID string, keyword string, page, pageSize int) ([]*models.Machine, int64, error)
	UpdateBusinessStatus(id string, status enums.BusinessStatus) error
	CheckDeviceExists(deviceID string) (bool, error)
}

// MachineRepository 售货机仓储实现
type MachineRepository struct {
	db *gorm.DB
}

// NewMachineRepository 创建售货机仓储
func NewMachineRepository(db *gorm.DB) MachineRepositoryInterface {
	return &MachineRepository{
		db: db,
	}
}

// GetByID 根据ID获取售货机
func (r *MachineRepository) GetByID(id string) (*models.Machine, error) {
	var machine models.Machine
	err := r.db.Where("id = ?", id).
		Preload("MachineOwner").
		First(&machine).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get machine by id: %w", err)
	}

	return &machine, nil
}

// GetByDeviceID 根据设备ID获取售货机
func (r *MachineRepository) GetByDeviceID(deviceID string) (*models.Machine, error) {
	var machine models.Machine
	err := r.db.Where("device_id = ?", deviceID).
		First(&machine).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get machine by device id: %w", err)
	}

	return &machine, nil
}

// GetList 获取售货机列表（根据机主ID）
func (r *MachineRepository) GetList(machineOwnerID string) ([]*models.Machine, error) {
	var machines []*models.Machine
	err := r.db.Where("machine_owner_id = ?", machineOwnerID).
		Order("created_at DESC").
		Find(&machines).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get machine list: %w", err)
	}

	return machines, nil
}

// GetPaging 分页获取售货机列表
func (r *MachineRepository) GetPaging(
	machineOwnerID string, keyword string, page, pageSize int,
) ([]*models.Machine, int64, error) {
	var machines []*models.Machine
	var totalCount int64

	query := r.db.Where("machine_owner_id = ?", machineOwnerID)

	// 添加关键词搜索
	if keyword != "" {
		keyword = strings.TrimSpace(keyword)
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR machine_no LIKE ? OR area LIKE ? OR address LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// 获取总数
	err := query.Model(&models.Machine{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count machines: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&machines).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get machines paging: %w", err)
	}

	return machines, totalCount, nil
}

// UpdateBusinessStatus 更新营业状态
func (r *MachineRepository) UpdateBusinessStatus(id string, status enums.BusinessStatus) error {
	result := r.db.Model(&models.Machine{}).
		Where("id = ?", id).
		Update("business_status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update business status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("machine not found")
	}

	return nil
}

// CheckDeviceExists 检查设备是否存在
func (r *MachineRepository) CheckDeviceExists(deviceID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Machine{}).
		Where("device_id = ?", deviceID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check device exists: %w", err)
	}

	return count > 0, nil
}
