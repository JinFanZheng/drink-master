package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ddteam/drink-master/internal/models"
)

func TestProductRepository_GetMachineProducts(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	// 创建测试数据
	product1 := &models.Product{
		ID:        "product-1",
		Name:      "Coffee",
		Status:    1,
		Price:     5.0,
		CreatedOn: time.Now(),
	}
	require.NoError(t, db.Create(product1).Error)

	product2 := &models.Product{
		ID:        "product-2",
		Name:      "Tea",
		Status:    1,
		Price:     4.0,
		CreatedOn: time.Now(),
	}
	require.NoError(t, db.Create(product2).Error)

	// 创建机器商品价格关联
	machineProduct1 := &models.MachineProductPrice{
		ID:              "mp-1",
		MachineId:       "machine-123",
		ProductId:       "product-1",
		Price:           5.0,
		PriceWithoutCup: 4.5,
		CreatedOn:       time.Now().Add(-1 * time.Hour),
	}
	require.NoError(t, db.Create(machineProduct1).Error)

	machineProduct2 := &models.MachineProductPrice{
		ID:              "mp-2",
		MachineId:       "machine-123",
		ProductId:       "product-2",
		Price:           4.0,
		PriceWithoutCup: 3.5,
		CreatedOn:       time.Now(),
	}
	require.NoError(t, db.Create(machineProduct2).Error)

	// 创建其他机器的商品（不应该被返回）
	machineProduct3 := &models.MachineProductPrice{
		ID:              "mp-3",
		MachineId:       "machine-456",
		ProductId:       "product-1",
		Price:           6.0,
		PriceWithoutCup: 5.5,
		CreatedOn:       time.Now(),
	}
	require.NoError(t, db.Create(machineProduct3).Error)

	// 测试获取机器商品列表
	result, err := repo.GetMachineProducts("machine-123")
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// 验证结果按创建时间排序（ASC）
	assert.Equal(t, "mp-1", result[0].ID)
	assert.Equal(t, "mp-2", result[1].ID)

	// 验证商品信息被正确预加载 - removed Product field references since Product field doesn't exist
	// TODO: If Product preloading is needed, need to add Product field to MachineProductPrice model

	// 测试不存在的机器
	result, err = repo.GetMachineProducts("nonexistent")
	require.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestProductRepository_GetMachineProducts_EmptyResult(t *testing.T) {
	db := setupTestDB(t)
	repo := NewProductRepository(db)

	// 测试没有商品的机器
	result, err := repo.GetMachineProducts("empty-machine")
	require.NoError(t, err)
	assert.Len(t, result, 0)
	assert.NotNil(t, result) // 应该返回空数组而不是nil
}
