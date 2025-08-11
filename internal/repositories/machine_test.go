package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 迁移表
	err = db.AutoMigrate(
		&models.MachineOwner{},
		&models.Machine{},
		&models.Product{},
		&models.MachineProductPrice{},
	)
	require.NoError(t, err)

	return db
}

func TestMachineRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMachineRepository(db)

	// 创建测试数据
	owner := &models.MachineOwner{
		ID:   "owner-123",
		Name: "Test Owner",
	}
	require.NoError(t, db.Create(owner).Error)

	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: "owner-123",
		MachineNo:      "M001",
		Name:           "Test Machine",
		Area:           "Test Area",
		Address:        "Test Address",
		BusinessStatus: "Open",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	require.NoError(t, db.Create(machine).Error)

	// 测试获取存在的机器
	result, err := repo.GetByID("machine-123")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "machine-123", result.ID)
	assert.Equal(t, "M001", result.MachineNo)
	assert.Equal(t, "Test Machine", result.Name)

	// 测试获取不存在的机器
	result, err = repo.GetByID("nonexistent")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestMachineRepository_GetByDeviceID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMachineRepository(db)

	deviceID := "device-123"
	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: "owner-123",
		MachineNo:      "M001",
		Name:           "Test Machine",
		DeviceId:       &deviceID,
		BusinessStatus: "Open",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	require.NoError(t, db.Create(machine).Error)

	// 测试获取存在的设备
	result, err := repo.GetByDeviceID("device-123")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "machine-123", result.ID)
	assert.Equal(t, "device-123", *result.DeviceId)

	// 测试获取不存在的设备
	result, err = repo.GetByDeviceID("nonexistent")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestMachineRepository_GetList(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMachineRepository(db)

	// 创建测试数据
	machines := []*models.Machine{
		{
			ID:             "machine-1",
			MachineOwnerId: "owner-123",
			MachineNo:      "M001",
			Name:           "Machine 1",
			BusinessStatus: "Open",
			CreatedAt:      time.Now().Add(-2 * time.Hour),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             "machine-2",
			MachineOwnerId: "owner-123",
			MachineNo:      "M002",
			Name:           "Machine 2",
			BusinessStatus: "Close",
			CreatedAt:      time.Now().Add(-1 * time.Hour),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             "machine-3",
			MachineOwnerId: "owner-456", // 不同的机主
			MachineNo:      "M003",
			Name:           "Machine 3",
			BusinessStatus: "Open",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for _, machine := range machines {
		require.NoError(t, db.Create(machine).Error)
	}

	// 测试获取特定机主的机器列表
	result, err := repo.GetList("owner-123")
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// 验证按创建时间倒序排列
	assert.Equal(t, "machine-2", result[0].ID) // 最新创建的
	assert.Equal(t, "machine-1", result[1].ID)

	// 测试不存在的机主
	result, err = repo.GetList("nonexistent")
	require.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestMachineRepository_GetPaging(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMachineRepository(db)

	// 创建测试数据
	machines := []*models.Machine{
		{
			ID:             "machine-1",
			MachineOwnerId: "owner-123",
			MachineNo:      "M001",
			Name:           "Coffee Machine",
			Area:           "Building A",
			Address:        "Floor 1",
			BusinessStatus: "Open",
			CreatedAt:      time.Now().Add(-3 * time.Hour),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             "machine-2",
			MachineOwnerId: "owner-123",
			MachineNo:      "M002",
			Name:           "Juice Machine",
			Area:           "Building B",
			Address:        "Floor 2",
			BusinessStatus: "Close",
			CreatedAt:      time.Now().Add(-2 * time.Hour),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             "machine-3",
			MachineOwnerId: "owner-123",
			MachineNo:      "M003",
			Name:           "Snack Machine",
			Area:           "Building A",
			Address:        "Floor 3",
			BusinessStatus: "Open",
			CreatedAt:      time.Now().Add(-1 * time.Hour),
			UpdatedAt:      time.Now(),
		},
	}

	for _, machine := range machines {
		require.NoError(t, db.Create(machine).Error)
	}

	// 测试不带搜索的分页
	result, totalCount, err := repo.GetPaging("owner-123", "", 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(3), totalCount)
	assert.Len(t, result, 2)
	assert.Equal(t, "machine-3", result[0].ID) // 最新创建的

	// 测试带搜索关键词的分页
	result, totalCount, err = repo.GetPaging("owner-123", "Coffee", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), totalCount)
	assert.Len(t, result, 1)
	assert.Equal(t, "Coffee Machine", result[0].Name)

	// 测试按区域搜索
	result, totalCount, err = repo.GetPaging("owner-123", "Building A", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), totalCount)
	assert.Len(t, result, 2)
}

func TestMachineRepository_UpdateBusinessStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMachineRepository(db)

	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: "owner-123",
		MachineNo:      "M001",
		Name:           "Test Machine",
		BusinessStatus: "Open",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	require.NoError(t, db.Create(machine).Error)

	// 测试更新状态
	err := repo.UpdateBusinessStatus("machine-123", "Close")
	require.NoError(t, err)

	// 验证更新结果
	var updated models.Machine
	require.NoError(t, db.First(&updated, "id = ?", "machine-123").Error)
	assert.Equal(t, "Close", updated.BusinessStatus)

	// 测试更新不存在的机器
	err = repo.UpdateBusinessStatus("nonexistent", "Open")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "machine not found")
}

func TestMachineRepository_CheckDeviceExists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMachineRepository(db)

	deviceID := "device-123"
	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: "owner-123",
		MachineNo:      "M001",
		Name:           "Test Machine",
		DeviceId:       &deviceID,
		BusinessStatus: "Open",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	require.NoError(t, db.Create(machine).Error)

	// 测试存在的设备
	exists, err := repo.CheckDeviceExists("device-123")
	require.NoError(t, err)
	assert.True(t, exists)

	// 测试不存在的设备
	exists, err = repo.CheckDeviceExists("nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)
}
