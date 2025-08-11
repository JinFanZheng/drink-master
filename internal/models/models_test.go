package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAutoMigrate(t *testing.T) {
	// 创建内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 测试自动迁移
	err = AutoMigrate(db)
	if err != nil {
		t.Errorf("AutoMigrate failed: %v", err)
	}

	// 验证表是否被创建
	tables := []string{"members", "machine_owners", "machines", "products", "machine_product_prices", "orders"}
	
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("Table %s was not created", table)
		}
	}
}

func TestAllModels(t *testing.T) {
	models := AllModels()
	
	if len(models) == 0 {
		t.Error("AllModels returned empty slice")
	}

	// 验证返回的模型数量
	expectedCount := 6 // Member, MachineOwner, Machine, Product, MachineProductPrice, Order
	if len(models) != expectedCount {
		t.Errorf("Expected %d models, got %d", expectedCount, len(models))
	}
}