package services

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = models.AutoMigrate(db)
	require.NoError(t, err)

	return db
}

func setupTestData(t *testing.T, db *gorm.DB) (string, string, string) {
	// 创建机主
	owner := models.MachineOwner{
		ID:               "owner-001",
		Name:             "Test Owner",
		ReceivingAccount: nil,
	}
	require.NoError(t, db.Create(&owner).Error)

	// 创建用户
	member := models.Member{
		ID:             "member-001",
		Nickname:       "Test User",
		WeChatOpenId:   "openid-001",
		Role:           "Member",
		MachineOwnerId: &owner.ID,
	}
	require.NoError(t, db.Create(&member).Error)

	// 创建机器
	machine := models.Machine{
		ID:             "machine-001",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM001",
		Name:           "Test Machine 1",
		Area:           "Test Area",
		Address:        "Test Address",
		BusinessStatus: "Open",
	}
	require.NoError(t, db.Create(&machine).Error)

	// 创建第二台机器
	machine2 := models.Machine{
		ID:             "machine-002",
		MachineOwnerId: owner.ID,
		MachineNo:      "VM002",
		Name:           "Test Machine 2",
		Area:           "Test Area 2",
		Address:        "Test Address 2",
		BusinessStatus: "Open",
	}
	require.NoError(t, db.Create(&machine2).Error)

	// 创建商品
	product := models.Product{
		ID:   "product-001",
		Name: "Test Coffee",
	}
	require.NoError(t, db.Create(&product).Error)

	return owner.ID, machine.ID, machine2.ID
}

func TestNewMachineOwnerService(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)

	assert.NotNil(t, service)
	assert.Equal(t, db, service.db)
}

func TestMachineOwnerService_GetSales(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)
	ownerID, machineID1, machineID2 := setupTestData(t, db)

	// 创建今天的订单数据
	today := time.Now().Truncate(24 * time.Hour)
	orders := []models.Order{
		{
			ID:            "order-001",
			MemberId:      "member-001",
			MachineId:     machineID1,
			ProductId:     "product-001",
			OrderNo:       "ON001",
			TotalAmount:   15.50,
			PayAmount:     15.50,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
		{
			ID:            "order-002",
			MemberId:      "member-001",
			MachineId:     machineID1,
			ProductId:     "product-001",
			OrderNo:       "ON002",
			TotalAmount:   12.00,
			PayAmount:     12.00,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
		{
			ID:            "order-003",
			MemberId:      "member-001",
			MachineId:     machineID2,
			ProductId:     "product-001",
			OrderNo:       "ON003",
			TotalAmount:   18.00,
			PayAmount:     18.00,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
	}

	for _, order := range orders {
		require.NoError(t, db.Create(&order).Error)
	}

	// 测试获取销售数据
	sales, err := service.GetSales(ownerID, today)
	require.NoError(t, err)
	require.Len(t, sales, 2) // 两台机器

	// 验证销售数据
	salesMap := make(map[string]decimal.Decimal)
	for _, sale := range sales {
		salesMap[sale.Label] = sale.Value
	}

	assert.True(t, salesMap["Test Machine 1"].Equal(decimal.NewFromFloat(27.50))) // 15.50 + 12.00
	assert.True(t, salesMap["Test Machine 2"].Equal(decimal.NewFromFloat(18.00)))
}

func TestMachineOwnerService_GetSales_EmptyOwnerID(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)

	sales, err := service.GetSales("", time.Now())
	assert.Error(t, err)
	assert.Nil(t, sales)
	assert.Contains(t, err.Error(), "机主ID不能为空")
}

func TestMachineOwnerService_GetSales_OwnerNotFound(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)

	sales, err := service.GetSales("nonexistent", time.Now())
	assert.Error(t, err)
	assert.Nil(t, sales)
	assert.Contains(t, err.Error(), "机主不存在")
}

func TestMachineOwnerService_GetSales_NoMachines(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)

	// 创建机主但没有机器
	owner := models.MachineOwner{
		ID:   "owner-empty",
		Name: "Empty Owner",
	}
	require.NoError(t, db.Create(&owner).Error)

	sales, err := service.GetSales(owner.ID, time.Now())
	require.NoError(t, err)
	assert.Empty(t, sales)
}

func TestMachineOwnerService_GetSales_NoSalesData(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)
	ownerID, _, _ := setupTestData(t, db)

	// 查询今天的销售数据 (没有订单)
	sales, err := service.GetSales(ownerID, time.Now())
	require.NoError(t, err)
	require.Len(t, sales, 2) // 两台机器

	// 验证销售额都是0
	for _, sale := range sales {
		assert.True(t, sale.Value.Equal(decimal.NewFromInt(0)))
	}
}

func TestMachineOwnerService_GetSales_OnlyPaidOrders(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)
	ownerID, machineID1, _ := setupTestData(t, db)

	today := time.Now().Truncate(24 * time.Hour)

	// 创建不同状态的订单
	orders := []models.Order{
		{
			ID:            "order-paid",
			MemberId:      "member-001",
			MachineId:     machineID1,
			ProductId:     "product-001",
			OrderNo:       "ON-PAID",
			PayAmount:     15.50,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
		{
			ID:            "order-unpaid",
			MemberId:      "member-001",
			MachineId:     machineID1,
			ProductId:     "product-001",
			OrderNo:       "ON-UNPAID",
			PayAmount:     12.00,
			PaymentStatus: "WaitPay",
			PaymentTime:   nil,
		},
		{
			ID:            "order-refunded",
			MemberId:      "member-001",
			MachineId:     machineID1,
			ProductId:     "product-001",
			OrderNo:       "ON-REFUND",
			PayAmount:     10.00,
			PaymentStatus: "Refunded",
			PaymentTime:   &today,
		},
	}

	for _, order := range orders {
		require.NoError(t, db.Create(&order).Error)
	}

	// 获取销售数据
	sales, err := service.GetSales(ownerID, today)
	require.NoError(t, err)

	// 查找机器1的销售数据
	var machine1Sales decimal.Decimal
	for _, sale := range sales {
		if sale.Label == "Test Machine 1" {
			machine1Sales = sale.Value
			break
		}
	}

	// 只有Paid状态的订单应该被计入销售额
	assert.True(t, machine1Sales.Equal(decimal.NewFromFloat(15.50)))
}

func TestMachineOwnerService_GetSalesStats(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)
	ownerID, machineID1, machineID2 := setupTestData(t, db)

	today := time.Now().Truncate(24 * time.Hour)

	// 创建订单
	orders := []models.Order{
		{
			ID:            "order-001",
			MemberId:      "member-001",
			MachineId:     machineID1,
			ProductId:     "product-001",
			OrderNo:       "ON001",
			PayAmount:     15.50,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
		{
			ID:            "order-002",
			MemberId:      "member-001",
			MachineId:     machineID2,
			ProductId:     "product-001",
			OrderNo:       "ON002",
			PayAmount:     12.00,
			PaymentStatus: "Paid",
			PaymentTime:   &today,
		},
	}

	for _, order := range orders {
		require.NoError(t, db.Create(&order).Error)
	}

	// 获取统计数据
	stats, err := service.GetSalesStats(ownerID, today, today)
	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, today, stats.Date)
	assert.Len(t, stats.Sales, 2)
	assert.True(t, stats.Total.Equal(decimal.NewFromFloat(27.50))) // 15.50 + 12.00
}

func TestMachineOwnerService_ValidateMachineOwnership(t *testing.T) {
	db := setupTestDB(t)
	service := NewMachineOwnerService(db)
	ownerID, machineID1, _ := setupTestData(t, db)

	// 测试有效的所有权
	err := service.ValidateMachineOwnership(ownerID, machineID1)
	assert.NoError(t, err)

	// 测试无效的所有权
	err = service.ValidateMachineOwnership("wrong-owner", machineID1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "您没有权限访问该机器")

	// 测试不存在的机器
	err = service.ValidateMachineOwnership(ownerID, "nonexistent-machine")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "您没有权限访问该机器")
}
