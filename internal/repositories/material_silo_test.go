package repositories

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

func TestMaterialSiloRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

	// Create test machine and product
	machine := &models.Machine{
		ID:             "test_machine",
		MachineOwnerId: "test_owner",
		MachineNo:      "M001",
		Name:           "Test Machine",
		BusinessStatus: enums.BusinessStatusOpen,
	}
	require.NoError(t, db.Create(machine).Error)

	product := &models.Product{
		ID:       "test_product",
		Name:     "Test Product",
		Category: func() *string { s := "drink"; return &s }(),
	}
	require.NoError(t, db.Create(product).Error)

	// Create test material silo
	silo := &models.MaterialSilo{
		ID:          "test_silo",
		MachineID:   machine.ID,
		SiloNo:      1,
		ProductID:   &product.ID,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  enums.SaleStatusOn,
	}
	require.NoError(t, db.Create(silo).Error)

	// Test GetByID - found
	result, err := repo.GetByID("test_silo")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test_silo", result.ID)
	assert.Equal(t, "test_machine", result.MachineID)
	assert.Equal(t, 1, result.SiloNo)
	assert.NotNil(t, result.ProductID)
	assert.Equal(t, "test_product", *result.ProductID)

	// Test GetByID - not found
	result, err = repo.GetByID("nonexistent")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestMaterialSiloRepository_GetByMachineID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

	// Create test machine
	machine := &models.Machine{
		ID:             "test_machine",
		MachineOwnerId: "test_owner",
		MachineNo:      "M001",
		Name:           "Test Machine",
		BusinessStatus: enums.BusinessStatusOpen,
	}
	require.NoError(t, db.Create(machine).Error)

	// Create multiple material silos
	silos := []*models.MaterialSilo{
		{
			ID:          "silo_1",
			MachineID:   machine.ID,
			SiloNo:      1,
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

	// Test GetByMachineID
	results, err := repo.GetByMachineID(machine.ID)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Should be ordered by silo_no ASC
	assert.Equal(t, 1, results[0].SiloNo)
	assert.Equal(t, 2, results[1].SiloNo)
}

func TestMaterialSiloRepository_GetPaging(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

	// Create test machine
	machine := &models.Machine{
		ID:             "test_machine",
		MachineOwnerId: "test_owner",
		MachineNo:      "M001",
		Name:           "Test Machine",
		BusinessStatus: enums.BusinessStatusOpen,
	}
	require.NoError(t, db.Create(machine).Error)

	// Create multiple material silos
	for i := 1; i <= 15; i++ {
		silo := &models.MaterialSilo{
			ID:          fmt.Sprintf("silo_%d", i),
			MachineID:   machine.ID,
			SiloNo:      i,
			Stock:       50,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOn,
		}
		require.NoError(t, db.Create(silo).Error)
	}

	// Test first page
	results, totalCount, err := repo.GetPaging(machine.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(15), totalCount)
	assert.Len(t, results, 10)

	// Test second page
	results, totalCount, err = repo.GetPaging(machine.ID, 2, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(15), totalCount)
	assert.Len(t, results, 5)
}

func TestMaterialSiloRepository_UpdateStock(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

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

	// Update stock
	err := repo.UpdateStock("test_silo", 75)
	require.NoError(t, err)

	// Verify update
	result, err := repo.GetByID("test_silo")
	require.NoError(t, err)
	assert.Equal(t, 75, result.Stock)
}

func TestMaterialSiloRepository_UpdateProduct(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

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

	// Update product
	err := repo.UpdateProduct("test_silo", "new_product")
	require.NoError(t, err)

	// Verify update
	result, err := repo.GetByID("test_silo")
	require.NoError(t, err)
	require.NotNil(t, result.ProductID)
	assert.Equal(t, "new_product", *result.ProductID)
}

func TestMaterialSiloRepository_UpdateSaleStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

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

	// Update sale status
	err := repo.UpdateSaleStatus("test_silo", enums.SaleStatusOff)
	require.NoError(t, err)

	// Verify update
	result, err := repo.GetByID("test_silo")
	require.NoError(t, err)
	assert.Equal(t, enums.SaleStatusOff, result.SaleStatus)
}

func TestMaterialSiloRepository_GetBySiloNo(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

	// Create test material silo
	silo := &models.MaterialSilo{
		ID:          "test_silo",
		MachineID:   "test_machine",
		SiloNo:      5,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  enums.SaleStatusOn,
	}
	require.NoError(t, db.Create(silo).Error)

	// Test GetBySiloNo - found
	result, err := repo.GetBySiloNo("test_machine", 5)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test_silo", result.ID)
	assert.Equal(t, 5, result.SiloNo)

	// Test GetBySiloNo - not found
	result, err = repo.GetBySiloNo("test_machine", 999)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestMaterialSiloRepository_GetByMachineAndProduct(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

	productID := "test_product"

	// Create multiple material silos with same product
	silos := []*models.MaterialSilo{
		{
			ID:          "silo_1",
			MachineID:   "test_machine",
			SiloNo:      1,
			ProductID:   &productID,
			Stock:       50,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOn,
		},
		{
			ID:          "silo_2",
			MachineID:   "test_machine",
			SiloNo:      3,
			ProductID:   &productID,
			Stock:       30,
			MaxCapacity: 100,
			SaleStatus:  enums.SaleStatusOff,
		},
	}

	for _, silo := range silos {
		require.NoError(t, db.Create(silo).Error)
	}

	// Test GetByMachineAndProduct
	results, err := repo.GetByMachineAndProduct("test_machine", productID)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Should be ordered by silo_no ASC
	assert.Equal(t, 1, results[0].SiloNo)
	assert.Equal(t, 3, results[1].SiloNo)
}

func TestMaterialSiloRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMaterialSiloRepository(db)

	silo := &models.MaterialSilo{
		ID:          "test_silo",
		MachineID:   "test_machine",
		SiloNo:      1,
		Stock:       50,
		MaxCapacity: 100,
		SaleStatus:  enums.SaleStatusOn,
	}

	err := repo.Create(silo)
	require.NoError(t, err)

	// Verify creation
	result, err := repo.GetByID("test_silo")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test_silo", result.ID)
}
