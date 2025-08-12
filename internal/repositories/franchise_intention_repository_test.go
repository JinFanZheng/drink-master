package repositories

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

func setupFranchiseTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func createTestFranchiseIntention(t *testing.T, db *gorm.DB, memberID string) *models.FranchiseIntention {
	intention := &models.FranchiseIntention{
		ID:               "test-franchise-1",
		MemberID:         memberID,
		ContactName:      "张三",
		ContactPhone:     "13800138000",
		IntendedLocation: "北京市朝阳区",
		Status:           "Pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := db.Create(intention).Error
	if err != nil {
		t.Fatalf("failed to create test franchise intention: %v", err)
	}

	return intention
}

func TestNewFranchiseIntentionRepository(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	if repo == nil {
		t.Fatal("expected repository to be created")
	}

	if repo.db != db {
		t.Error("expected repository to have correct database connection")
	}
}

func TestFranchiseIntentionRepository_Create(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	intention := &models.FranchiseIntention{
		ID:               "test-create-1",
		MemberID:         "member-123",
		ContactName:      "李四",
		ContactPhone:     "13900139000",
		IntendedLocation: "上海市",
		Status:           "Pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := repo.Create(intention)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证创建
	var found models.FranchiseIntention
	err = db.Where("id = ?", intention.ID).First(&found).Error
	if err != nil {
		t.Fatalf("expected to find created intention, got error: %v", err)
	}

	if found.ContactName != intention.ContactName {
		t.Errorf("expected contact name '%s', got '%s'", intention.ContactName, found.ContactName)
	}
}

func TestFranchiseIntentionRepository_GetByID(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	// 创建测试数据
	testIntention := createTestFranchiseIntention(t, db, "member-123")

	// 测试正确的ID
	intention, err := repo.GetByID(testIntention.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if intention.ID != testIntention.ID {
		t.Errorf("expected ID '%s', got '%s'", testIntention.ID, intention.ID)
	}

	if intention.ContactName != testIntention.ContactName {
		t.Errorf("expected contact name '%s', got '%s'", testIntention.ContactName, intention.ContactName)
	}

	// 测试不存在的ID
	_, err = repo.GetByID("non-existent-id")
	if err == nil {
		t.Error("expected error for non-existent intention")
	}
}

func TestFranchiseIntentionRepository_GetByMemberID(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	memberID := "member-456"

	// 创建多个测试数据
	createTestFranchiseIntention(t, db, memberID)

	intention2 := &models.FranchiseIntention{
		ID:               "test-franchise-2",
		MemberID:         memberID,
		ContactName:      "王五",
		ContactPhone:     "13700137000",
		IntendedLocation: "深圳市",
		Status:           "Approved",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := repo.Create(intention2)
	if err != nil {
		t.Fatalf("failed to create second intention: %v", err)
	}

	// 测试获取会员的意向列表
	intentions, err := repo.GetByMemberID(memberID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(intentions) != 2 {
		t.Errorf("expected 2 intentions, got %d", len(intentions))
	}

	// 测试不存在的会员ID
	intentions, err = repo.GetByMemberID("non-existent-member")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(intentions) != 0 {
		t.Errorf("expected 0 intentions, got %d", len(intentions))
	}
}

func TestFranchiseIntentionRepository_CheckExistingByMember(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	memberID := "member-789"

	// 测试没有待处理意向的情况
	exists, err := repo.CheckExistingByMember(memberID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if exists {
		t.Error("expected no existing intention")
	}

	// 创建待处理意向
	createTestFranchiseIntention(t, db, memberID)

	// 测试有待处理意向的情况
	exists, err = repo.CheckExistingByMember(memberID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !exists {
		t.Error("expected existing intention")
	}
}

func TestFranchiseIntentionRepository_UpdateStatus(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	// 创建测试数据
	testIntention := createTestFranchiseIntention(t, db, "member-999")

	// 更新状态
	err := repo.UpdateStatus(testIntention.ID, "Approved")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证更新
	intention, err := repo.GetByID(testIntention.ID)
	if err != nil {
		t.Fatalf("expected no error when getting updated intention, got %v", err)
	}

	if intention.Status != "Approved" {
		t.Errorf("expected status 'Approved', got '%s'", intention.Status)
	}
}

func TestFranchiseIntentionRepository_Update(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	// 创建测试数据
	testIntention := createTestFranchiseIntention(t, db, "member-111")

	// 更新数据
	testIntention.ContactName = "更新的姓名"
	testIntention.ContactPhone = "13600136000"

	err := repo.Update(testIntention)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证更新
	intention, err := repo.GetByID(testIntention.ID)
	if err != nil {
		t.Fatalf("expected no error when getting updated intention, got %v", err)
	}

	if intention.ContactName != "更新的姓名" {
		t.Errorf("expected contact name '更新的姓名', got '%s'", intention.ContactName)
	}
}

func TestFranchiseIntentionRepository_Delete(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	// 创建测试数据
	testIntention := createTestFranchiseIntention(t, db, "member-222")

	// 删除数据
	err := repo.Delete(testIntention.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证删除
	_, err = repo.GetByID(testIntention.ID)
	if err == nil {
		t.Error("expected error when getting deleted intention")
	}
}
