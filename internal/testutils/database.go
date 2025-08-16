package testutils

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ddteam/drink-master/internal/models"
)

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// SetupTestDB 创建测试用的SQLite内存数据库
func SetupTestDB(t *testing.T) *gorm.DB {
	// 创建SQLite内存数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 静默模式，减少测试输出
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 执行数据库迁移
	if err := models.AutoMigrate(db); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// TeardownTestDB 清理测试数据库（SQLite内存数据库会自动清理）
func TeardownTestDB(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// SeedTestData 为测试数据库填充基础测试数据
func SeedTestData(t *testing.T, db *gorm.DB) {
	// 创建测试用户
	testMember := &models.Member{
		ID:           "test-member-001",
		Nickname:     stringPtr("Test User"),
		Avatar:       stringPtr("https://example.com/avatar.jpg"),
		WeChatOpenId: stringPtr("test-openid-123"),
		Role:         1,
		IsAdmin:      models.BitBool(0),
		CreatedOn:    time.Now(),
	}
	if err := db.Create(testMember).Error; err != nil {
		t.Logf("Warning: Failed to create test member: %v", err)
	}

	// 创建测试机器
	testMachine := &models.Machine{
		ID:             "test-machine-001",
		Name:           stringPtr("测试饮品机"),
		Area:           stringPtr("测试区域"),
		Address:        stringPtr("测试地点"),
		BusinessStatus: 1, // 营业中
		CreatedOn:      time.Now(),
	}
	if err := db.Create(testMachine).Error; err != nil {
		t.Logf("Warning: Failed to create test machine: %v", err)
	}

	// 创建测试产品
	testProduct := &models.Product{
		ID:              "test-product-001",
		Name:            "测试饮品",
		Status:          1,
		Price:           5.50,
		PriceWithoutCup: 5.00,
		CreatedOn:       time.Now(),
	}
	if err := db.Create(testProduct).Error; err != nil {
		t.Logf("Warning: Failed to create test product: %v", err)
	}

	// 创建机器产品价格关联
	testPrice := &models.MachineProductPrice{
		ID:              "test-price-001",
		MachineId:       "test-machine-001",
		ProductId:       "test-product-001",
		Price:           5.50,
		PriceWithoutCup: 5.00,
		CreatedOn:       time.Now(),
	}
	if err := db.Create(testPrice).Error; err != nil {
		t.Logf("Warning: Failed to create test machine product price: %v", err)
	}
}
