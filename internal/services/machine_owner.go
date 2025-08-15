package services

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

// MachineOwnerService 机主服务
type MachineOwnerService struct {
	db *gorm.DB
}

// NewMachineOwnerService 创建机主服务
func NewMachineOwnerService(db *gorm.DB) *MachineOwnerService {
	return &MachineOwnerService{
		db: db,
	}
}

// GetSales 获取机主的销售情况数据
func (s *MachineOwnerService) GetSales(machineOwnerID string, targetDate time.Time) ([]contracts.ColumnModel, error) {
	if machineOwnerID == "" {
		return nil, fmt.Errorf("机主ID不能为空")
	}

	// 验证机主是否存在
	var machineOwner models.MachineOwner
	if err := s.db.Where("id = ?", machineOwnerID).First(&machineOwner).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("机主不存在")
		}
		return nil, fmt.Errorf("查询机主信息失败: %w", err)
	}

	// 获取机主的所有机器
	var machines []models.Machine
	if err := s.db.Where("machine_owner_id = ? AND deleted_at IS NULL", machineOwnerID).Find(&machines).Error; err != nil {
		return nil, fmt.Errorf("查询机器列表失败: %w", err)
	}

	if len(machines) == 0 {
		return []contracts.ColumnModel{}, nil
	}

	// 构建机器ID列表
	machineIDs := make([]string, len(machines))
	machineMap := make(map[string]string) // machine_id -> machine_name
	for i, machine := range machines {
		machineIDs[i] = machine.ID
		if machine.Name != nil {
			machineMap[machine.ID] = *machine.Name
		} else {
			machineMap[machine.ID] = ""
		}
	}

	// 设置日期范围 (整天)
	startDate := targetDate.Truncate(24 * time.Hour)
	endDate := startDate.Add(24 * time.Hour)

	// 查询当天的订单销售数据
	type SalesData struct {
		MachineID  string          `db:"machine_id"`
		TotalSales decimal.Decimal `db:"total_sales"`
		OrderCount int64           `db:"order_count"`
	}

	var salesData []SalesData
	err := s.db.Table("orders").
		Select("machine_id, SUM(pay_amount) as total_sales, COUNT(*) as order_count").
		Where("machine_id IN ? AND payment_status = ? AND payment_time >= ? AND payment_time < ? AND deleted_at IS NULL",
			machineIDs, int(enums.PaymentStatusPaid), startDate, endDate).
		Group("machine_id").
		Scan(&salesData).Error

	if err != nil {
		return nil, fmt.Errorf("查询销售数据失败: %w", err)
	}

	// 构建销售数据映射
	salesMap := make(map[string]decimal.Decimal)
	for _, data := range salesData {
		salesMap[data.MachineID] = data.TotalSales
	}

	// 构建返回结果 (包含所有机器，即使销售额为0)
	result := make([]contracts.ColumnModel, 0, len(machines))
	for _, machine := range machines {
		sales := salesMap[machine.ID]
		if sales.IsZero() {
			sales = decimal.NewFromInt(0)
		}

		label := ""
		if machine.Name != nil {
			label = *machine.Name
		}
		result = append(result, contracts.ColumnModel{
			Label: label,
			Value: sales,
		})
	}

	return result, nil
}

// GetSalesStats 获取机主销售统计 (可扩展功能)
func (s *MachineOwnerService) GetSalesStats(
	machineOwnerID string,
	startDate, endDate time.Time,
) (*contracts.SalesResponse, error) {
	sales, err := s.GetSales(machineOwnerID, startDate)
	if err != nil {
		return nil, err
	}

	// 计算总销售额
	total := decimal.NewFromInt(0)
	for _, sale := range sales {
		total = total.Add(sale.Value)
	}

	return &contracts.SalesResponse{
		Date:  startDate,
		Sales: sales,
		Total: total,
	}, nil
}

// ValidateMachineOwnership 验证机器所有权
func (s *MachineOwnerService) ValidateMachineOwnership(machineOwnerID, machineID string) error {
	var count int64
	err := s.db.Model(&models.Machine{}).
		Where("id = ? AND machine_owner_id = ? AND deleted_at IS NULL", machineID, machineOwnerID).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("验证机器所有权失败: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("您没有权限访问该机器")
	}

	return nil
}
