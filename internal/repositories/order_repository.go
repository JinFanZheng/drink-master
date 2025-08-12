package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

// OrderRepository 订单仓库接口
type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id string) (*models.Order, error)
	GetByMemberPaging(memberID string, pageIndex, pageSize int) ([]models.Order, int64, error)
	Update(order *models.Order) error
	Delete(id string) error
	GetByOrderNo(orderNo string) (*models.Order, error)
}

// orderRepository 订单仓库实现
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓库
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// Create 创建订单
func (r *orderRepository) Create(order *models.Order) error {
	if order.ID == "" {
		order.ID = uuid.New().String()
	}
	return r.db.Create(order).Error
}

// GetByID 根据ID获取订单
func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Member").Preload("Machine").Preload("Product").
		Where("id = ? AND deleted_at IS NULL", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByMemberPaging 分页获取会员订单列表
func (r *orderRepository) GetByMemberPaging(memberID string, pageIndex, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// 计算总数
	err := r.db.Model(&models.Order{}).
		Where("member_id = ? AND deleted_at IS NULL", memberID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (pageIndex - 1) * pageSize

	// 查询订单列表
	err = r.db.Preload("Product").
		Where("member_id = ? AND deleted_at IS NULL", memberID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&orders).Error

	return orders, total, err
}

// Update 更新订单
func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// Delete 软删除订单
func (r *orderRepository) Delete(id string) error {
	return r.db.Model(&models.Order{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

// GetByOrderNo 根据订单号获取订单
func (r *orderRepository) GetByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Member").Preload("Machine").Preload("Product").
		Where("order_no = ? AND deleted_at IS NULL", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
