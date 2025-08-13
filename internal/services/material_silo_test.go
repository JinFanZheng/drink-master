package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

func TestMaterialSiloService_GetPaging(t *testing.T) {
	db := setupTestDB(t)
	service := NewMaterialSiloService(db)

	// Create test machine
	machine := &models.Machine{
		ID:             "test_machine",
		MachineOwnerId: "test_owner",
		MachineNo:      "M001",
		Name:           "Test Machine",
		BusinessStatus: enums.BusinessStatusOpen,
	}
	require.NoError(t, db.Create(machine).Error)

	// Create test product
	product := &models.Product{
		ID:       "test_product",
		Name:     "Test Product",
		Category: func() *string { s := "drink"; return &s }(),
	}
	require.NoError(t, db.Create(product).Error)

	// Create test material silos
	silos := []*models.MaterialSilo{
		{
			ID:          "silo_1",
			MachineID:   machine.ID,
			SiloNo:      1,
			ProductID:   &product.ID,
			Stock:       50,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOn,
		},
		{
			ID:          "silo_2",
			MachineID:   machine.ID,
			SiloNo:      2,
			Stock:       30,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOff,
		},
	}

	for _, silo := range silos {
		require.NoError(t, db.Create(silo).Error)
	}

	// Test successful paging request
	req := contracts.GetMaterialSiloPagingRequest{
		MachineID: machine.ID,
		PageIndex: 1,
		PageSize:  10,
	}

	result, err := service.GetPaging(req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Len(t, result.Items, 2)
	assert.Equal(t, int64(2), result.TotalCount)
	assert.Equal(t, 1, result.PageIndex)
	assert.Equal(t, 10, result.PageSize)

	// Check first item details
	item1 := result.Items[0]
	assert.Equal(t, "silo_1", item1.ID)
	assert.Equal(t, machine.ID, item1.MachineID)
	assert.Equal(t, 1, item1.SiloNo)
	assert.NotNil(t, item1.ProductID)
	assert.Equal(t, product.ID, *item1.ProductID)
	assert.NotNil(t, item1.ProductName)
	assert.Equal(t, product.Name, *item1.ProductName)
	assert.Equal(t, 50, item1.Stock)
	assert.Equal(t, 100, item1.MaxCapacity)
	assert.Equal(t, "On", item1.SaleStatus)

	// Test with non-existent machine
	req.MachineID = "nonexistent"
	_, err = service.GetPaging(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "机器不存在")
}

func TestMaterialSiloService_UpdateStock(t *testing.T) {
	db := setupTestDB(t)
	service := NewMaterialSiloService(db)

	// Create test material silo
	silo := &models.MaterialSilo{
		ID:          "test_silo",
		MachineID:   "test_machine",
		SiloNo:      1,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  enums.SaleStatusOn,
	}
	require.NoError(t, db.Create(silo).Error)

	// Test successful stock update
	req := contracts.UpdateMaterialSiloStockRequest{
		ID:    "test_silo",
		Stock: 75,
	}

	result, err := service.UpdateStock(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "库存更新成功", result.Message)

	// Test stock exceeding capacity
	req.Stock = 150
	result, err = service.UpdateStock(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.Success)
	assert.Contains(t, result.Message, "库存不能超过最大容量")

	// Test with non-existent silo
	req.ID = "nonexistent"
	req.Stock = 50
	result, err = service.UpdateStock(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.Success)
	assert.Equal(t, "物料槽不存在", result.Message)
}

func TestMaterialSiloService_UpdateProduct(t *testing.T) {
	db := setupTestDB(t)
	service := NewMaterialSiloService(db)

	// Create test product
	product := &models.Product{
		ID:       "test_product",
		Name:     "Test Product",
		Category: func() *string { s := "drink"; return &s }(),
	}
	require.NoError(t, db.Create(product).Error)

	// Create test material silo
	silo := &models.MaterialSilo{
		ID:          "test_silo",
		MachineID:   "test_machine",
		SiloNo:      1,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  enums.SaleStatusOn,
	}
	require.NoError(t, db.Create(silo).Error)

	// Test successful product update
	req := contracts.UpdateMaterialSiloProductRequest{
		ID:        "test_silo",
		ProductID: product.ID,
	}

	result, err := service.UpdateProduct(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "产品更新成功", result.Message)

	// Test with non-existent product
	req.ProductID = "nonexistent"
	result, err = service.UpdateProduct(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.Success)
	assert.Equal(t, "产品不存在", result.Message)

	// Test with non-existent silo
	req.ID = "nonexistent"
	req.ProductID = product.ID
	result, err = service.UpdateProduct(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.Success)
	assert.Equal(t, "物料槽不存在", result.Message)
}

func TestMaterialSiloService_ToggleSaleStatus(t *testing.T) {
	db := setupTestDB(t)
	service := NewMaterialSiloService(db)

	// Create test product
	product := &models.Product{
		ID:       "test_product",
		Name:     "Test Product",
		Category: func() *string { s := "drink"; return &s }(),
	}
	require.NoError(t, db.Create(product).Error)

	tests := []struct {
		name           string
		silo           *models.MaterialSilo
		request        contracts.ToggleSaleMaterialSiloRequest
		expectedResult *contracts.MaterialSiloOperationResult
	}{
		{
			name: "turn off sale status",
			silo: &models.MaterialSilo{
				ID:          "test_silo_1",
				MachineID:   "test_machine",
				SiloNo:      1,
				ProductID:   &product.ID,
				Stock:       50,
				MaxCapacity: 100,
				SaleStatus:  enums.SaleStatusOn,
			},
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "test_silo_1",
				SaleStatus: contracts.SaleStatusOff,
			},
			expectedResult: &contracts.MaterialSiloOperationResult{
				Success: true,
				Message: "销售状态已切换为停售",
				Data:    "停售",
			},
		},
		{
			name: "turn on sale status with product and stock",
			silo: &models.MaterialSilo{
				ID:          "test_silo_2",
				MachineID:   "test_machine",
				SiloNo:      2,
				ProductID:   &product.ID,
				Stock:       50,
				MaxCapacity: 100,
				SaleStatus:  enums.SaleStatusOff,
			},
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "test_silo_2",
				SaleStatus: contracts.SaleStatusOn,
			},
			expectedResult: &contracts.MaterialSiloOperationResult{
				Success: true,
				Message: "销售状态已切换为在售",
				Data:    "在售",
			},
		},
		{
			name: "turn on sale status without product",
			silo: &models.MaterialSilo{
				ID:          "test_silo_3",
				MachineID:   "test_machine",
				SiloNo:      3,
				ProductID:   nil,
				Stock:       50,
				MaxCapacity: 100,
				SaleStatus:  enums.SaleStatusOff,
			},
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "test_silo_3",
				SaleStatus: contracts.SaleStatusOn,
			},
			expectedResult: &contracts.MaterialSiloOperationResult{
				Success: false,
				Message: "开启销售前需要先设置产品",
			},
		},
		{
			name: "turn on sale status without stock",
			silo: &models.MaterialSilo{
				ID:          "test_silo_4",
				MachineID:   "test_machine",
				SiloNo:      4,
				ProductID:   &product.ID,
				Stock:       0,
				MaxCapacity: 100,
				SaleStatus:  enums.SaleStatusOff,
			},
			request: contracts.ToggleSaleMaterialSiloRequest{
				ID:         "test_silo_4",
				SaleStatus: contracts.SaleStatusOn,
			},
			expectedResult: &contracts.MaterialSiloOperationResult{
				Success: false,
				Message: "开启销售前需要先补充库存",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test silo
			require.NoError(t, db.Create(tt.silo).Error)

			// Execute test
			result, err := service.ToggleSaleStatus(tt.request)
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tt.expectedResult.Success, result.Success)
			assert.Equal(t, tt.expectedResult.Message, result.Message)
			if tt.expectedResult.Data != "" {
				assert.Equal(t, tt.expectedResult.Data, result.Data)
			}
		})
	}

	// Test with non-existent silo
	req := contracts.ToggleSaleMaterialSiloRequest{
		ID:         "nonexistent",
		SaleStatus: contracts.SaleStatusOn,
	}
	result, err := service.ToggleSaleStatus(req)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.Success)
	assert.Equal(t, "物料槽不存在", result.Message)
}

func TestMaterialSiloService_ValidateMachineExists(t *testing.T) {
	db := setupTestDB(t)
	service := NewMaterialSiloService(db)

	// Create test machine
	machine := &models.Machine{
		ID:             "test_machine",
		MachineOwnerId: "test_owner",
		MachineNo:      "M001",
		Name:           "Test Machine",
		BusinessStatus: enums.BusinessStatusOpen,
	}
	require.NoError(t, db.Create(machine).Error)

	// Test existing machine
	err := service.ValidateMachineExists("test_machine")
	assert.NoError(t, err)

	// Test non-existent machine
	err = service.ValidateMachineExists("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "机器不存在")
}

func TestMaterialSiloService_ValidateProductExists(t *testing.T) {
	db := setupTestDB(t)
	service := NewMaterialSiloService(db)

	// Create test product
	product := &models.Product{
		ID:       "test_product",
		Name:     "Test Product",
		Category: func() *string { s := "drink"; return &s }(),
	}
	require.NoError(t, db.Create(product).Error)

	// Test existing product
	err := service.ValidateProductExists("test_product")
	assert.NoError(t, err)

	// Test non-existent product
	err = service.ValidateProductExists("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "产品不存在")
}
