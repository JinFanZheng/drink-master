package services

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

func TestNewMachineService_RealDB(t *testing.T) {
	// Create a real in-memory database for the constructor test
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	service := NewMachineService(db)

	if service == nil {
		t.Fatal("expected service to be created")
	}

	// Verify the service can be cast to the concrete type
	concreteService, ok := service.(*MachineService)
	if !ok {
		t.Fatal("expected service to be of type *MachineService")
	}

	if concreteService.machineRepo == nil {
		t.Error("expected machineRepo to be set")
	}

	if concreteService.productRepo == nil {
		t.Error("expected productRepo to be set")
	}

	if concreteService.deviceService == nil {
		t.Error("expected deviceService to be set")
	}
}
