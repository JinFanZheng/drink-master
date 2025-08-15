package repositories

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)


func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 迁移数据库
	err = models.AutoMigrate(db)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func createTestMember(t *testing.T, db *gorm.DB) *models.Member {
	member := &models.Member{
		ID:           "test-member-1",
		Nickname:     stringPtr("测试用户"),
		Avatar:       stringPtr("https://example.com/avatar.jpg"),
		WeChatOpenId: stringPtr("test-openid-1"),
		Role:         1, // Member role as int
		IsAdmin:      models.BitBool(0), // false as BitBool
		CreatedOn:    time.Now(),
	}

	err := db.Create(member).Error
	if err != nil {
		t.Fatalf("failed to create test member: %v", err)
	}

	return member
}

func TestNewMemberRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMemberRepository(db)

	if repo == nil {
		t.Fatal("expected repository to be created")
	}

	if repo.db != db {
		t.Error("expected repository to have correct database connection")
	}
}

func TestMemberRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMemberRepository(db)

	// 创建测试数据
	testMember := createTestMember(t, db)

	// 测试正确的ID
	member, err := repo.GetByID(testMember.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if member.ID != testMember.ID {
		t.Errorf("expected member ID '%s', got '%s'", testMember.ID, member.ID)
	}

	if member.Nickname != testMember.Nickname {
		t.Errorf("expected nickname '%s', got '%s'", testMember.Nickname, member.Nickname)
	}

	// 测试不存在的ID
	_, err = repo.GetByID("non-existent-id")
	if err == nil {
		t.Error("expected error for non-existent member")
	}
}

func TestMemberRepository_GetByWeChatOpenID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMemberRepository(db)

	// 创建测试数据
	testMember := createTestMember(t, db)

	// 测试正确的OpenID  
	member, err := repo.GetByWeChatOpenID(*testMember.WeChatOpenId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if member.WeChatOpenId != testMember.WeChatOpenId {
		t.Errorf("expected openID '%s', got '%s'", testMember.WeChatOpenId, member.WeChatOpenId)
	}

	// 测试不存在的OpenID
	_, err = repo.GetByWeChatOpenID("non-existent-openid")
	if err == nil {
		t.Error("expected error for non-existent member")
	}
}

func TestMemberRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMemberRepository(db)

	// 创建测试数据
	testMember := createTestMember(t, db)

	// 更新数据
	testMember.Nickname = stringPtr("更新的昵称")
	testMember.Avatar = stringPtr("https://example.com/new-avatar.jpg")

	err := repo.Update(testMember)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证更新
	updatedMember, err := repo.GetByID(testMember.ID)
	if err != nil {
		t.Fatalf("expected no error when getting updated member, got %v", err)
	}

	if *updatedMember.Nickname != "更新的昵称" {
		t.Errorf("expected updated nickname '更新的昵称', got '%s'", updatedMember.Nickname)
	}
}

func TestMemberRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewMemberRepository(db)

	member := &models.Member{
		ID:           "new-member-1",
		Nickname:     stringPtr("新用户"),
		Avatar:       stringPtr("https://example.com/avatar.jpg"),
		WeChatOpenId: stringPtr("new-openid-1"),
		Role:         1,
		IsAdmin:      models.BitBool(0),
		CreatedOn:    time.Now(),
	}

	err := repo.Create(member)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证创建
	createdMember, err := repo.GetByID(member.ID)
	if err != nil {
		t.Fatalf("expected no error when getting created member, got %v", err)
	}

	if createdMember.Nickname != member.Nickname {
		t.Errorf("expected nickname '%s', got '%s'", member.Nickname, createdMember.Nickname)
	}
}

func TestMemberRepository_GetMemberWithFranchiseIntentions(t *testing.T) {
	db := setupTestDB(t)
	memberRepo := NewMemberRepository(db)
	franchiseRepo := NewFranchiseIntentionRepository(db)

	// 创建测试会员
	testMember := createTestMember(t, db)

	// 创建加盟意向
	intention := &models.FranchiseIntention{
		ID:        "test-intention-1",
		MemberId:  stringPtr(testMember.ID),
		Name:      stringPtr("张三"),
		Mobile:    stringPtr("13800138000"),
		Area:      stringPtr("北京市"),
		CreatedOn: time.Now(),
	}

	err := franchiseRepo.Create(intention)
	if err != nil {
		t.Fatalf("failed to create test franchise intention: %v", err)
	}

	// 测试获取会员和加盟意向
	member, intentions, err := memberRepo.GetMemberWithFranchiseIntentions(testMember.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if member.ID != testMember.ID {
		t.Errorf("expected member ID '%s', got '%s'", testMember.ID, member.ID)
	}

	if len(intentions) != 1 {
		t.Errorf("expected 1 franchise intention, got %d", len(intentions))
	}

	if intentions[0].ID != intention.ID {
		t.Errorf("expected intention ID '%s', got '%s'", intention.ID, intentions[0].ID)
	}
}

func TestMemberRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	memberRepo := NewMemberRepository(db)

	// 创建测试会员
	testMember := createTestMember(t, db)

	// 删除会员
	err := memberRepo.Delete(testMember.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证删除
	_, err = memberRepo.GetByID(testMember.ID)
	if err == nil {
		t.Error("expected error when getting deleted member")
	}
}
