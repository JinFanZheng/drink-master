package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/enums"
	"github.com/ddteam/drink-master/internal/models"
)

func setupMachineTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = models.AutoMigrate(db)
	require.NoError(t, err)

	return db
}

func TestMachineRepository_GetByID(t *testing.T) {
	db := setupMachineTestDB(t)
	repo := NewMachineRepository(db)

	// Create test data
	ownerName := "Test Owner"
	owner := &models.MachineOwner{
		ID:   "owner-123",
		Name: &ownerName,
	}
	require.NoError(t, db.Create(owner).Error)

	machineNo := "M001"
	machineName := "Test Machine"
	area := "Test Area"
	address := "Test Address"
	now := time.Now()

	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: &owner.ID,
		MachineNo:      &machineNo,
		Name:           &machineName,
		Area:           &area,
		Address:        &address,
		BusinessStatus: enums.BusinessStatusOpen,
		CreatedOn:      now,
		UpdatedOn:      &now,
	}
	require.NoError(t, db.Create(machine).Error)

	// Test GetByID
	result, err := repo.GetByID("machine-123")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "machine-123", result.ID)
	assert.NotNil(t, result.MachineNo)
	assert.Equal(t, "M001", *result.MachineNo)
}

func TestMachineRepository_GetByOwnerID(t *testing.T) {
	db := setupMachineTestDB(t)
	repo := NewMachineRepository(db)

	// Create test data
	ownerName := "Test Owner"
	owner := &models.MachineOwner{
		ID:   "owner-123",
		Name: &ownerName,
	}
	require.NoError(t, db.Create(owner).Error)

	machineNo := "M001"
	machineName := "Test Machine"
	now := time.Now()

	machine := &models.Machine{
		ID:             "machine-123",
		MachineOwnerId: &owner.ID,
		MachineNo:      &machineNo,
		Name:           &machineName,
		BusinessStatus: enums.BusinessStatusOpen,
		CreatedOn:      now,
	}
	require.NoError(t, db.Create(machine).Error)

	// Test GetByOwnerID
	results, err := repo.GetByOwnerID("owner-123")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "machine-123", results[0].ID)
}

func TestMachineRepository_Update(t *testing.T) {
	db := setupMachineTestDB(t)
	repo := NewMachineRepository(db)

	// Create test data
	machineNo := "M001"
	machineName := "Test Machine"
	now := time.Now()

	machine := &models.Machine{
		ID:             "machine-123",
		MachineNo:      &machineNo,
		Name:           &machineName,
		BusinessStatus: enums.BusinessStatusOpen,
		CreatedOn:      now,
	}
	require.NoError(t, db.Create(machine).Error)

	// Update machine
	newName := "Updated Machine"
	machine.Name = &newName
	err := repo.Update(machine)
	assert.NoError(t, err)

	// Verify update
	result, err := repo.GetByID("machine-123")
	assert.NoError(t, err)
	assert.NotNil(t, result.Name)
	assert.Equal(t, "Updated Machine", *result.Name)
}