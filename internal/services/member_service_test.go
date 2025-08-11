package services

import (
	"testing"
	"time"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/models"
	"github.com/ddteam/drink-master/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupServiceTestDB(t *testing.T) *gorm.DB {
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

func createServiceTestMember(t *testing.T, db *gorm.DB) *models.Member {
	member := &models.Member{
		ID:           "service-test-member-1",
		Nickname:     "服务测试用户",
		Avatar:       "https://example.com/old-avatar.jpg",
		WeChatOpenId: "service-test-openid-1",
		Role:         "Member",
		IsAdmin:      false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := db.Create(member).Error
	if err != nil {
		t.Fatalf("failed to create test member: %v", err)
	}

	return member
}

func TestNewMemberService(t *testing.T) {
	db := setupServiceTestDB(t)
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)

	service := NewMemberService(memberRepo, franchiseRepo)

	if service == nil {
		t.Fatal("expected service to be created")
	}

	if service.memberRepo != memberRepo {
		t.Error("expected service to have correct member repository")
	}

	if service.franchiseIntentionRepo != franchiseRepo {
		t.Error("expected service to have correct franchise repository")
	}
}

func TestMemberService_UpdateMember(t *testing.T) {
	db := setupServiceTestDB(t)
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)
	service := NewMemberService(memberRepo, franchiseRepo)

	// 创建测试数据
	testMember := createServiceTestMember(t, db)

	// 测试正常更新
	req := contracts.UpdateMemberRequest{
		Nickname: "新昵称",
		Avatar:   "https://example.com/new-avatar.jpg",
	}

	response, err := service.UpdateMember(testMember.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Nickname != req.Nickname {
		t.Errorf("expected nickname '%s', got '%s'", req.Nickname, response.Nickname)
	}

	if response.Avatar != req.Avatar {
		t.Errorf("expected avatar '%s', got '%s'", req.Avatar, response.Avatar)
	}

	// 测试空昵称
	emptyReq := contracts.UpdateMemberRequest{
		Nickname: "",
		Avatar:   "https://example.com/avatar.jpg",
	}

	_, err = service.UpdateMember(testMember.ID, emptyReq)
	if err == nil {
		t.Error("expected error for empty nickname")
	}

	// 测试不存在的用户
	_, err = service.UpdateMember("non-existent-id", req)
	if err == nil {
		t.Error("expected error for non-existent member")
	}
}

func TestMemberService_CreateFranchiseIntention(t *testing.T) {
	db := setupServiceTestDB(t)
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)
	service := NewMemberService(memberRepo, franchiseRepo)

	// 创建测试数据
	testMember := createServiceTestMember(t, db)

	// 测试正常创建
	req := contracts.CreateFranchiseIntentionRequest{
		ContactName:      "张三",
		ContactPhone:     "13800138000",
		IntendedLocation: "北京市朝阳区",
		Remarks:          "希望在繁华地段开店",
	}

	response, err := service.CreateFranchiseIntention(testMember.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.ContactName != req.ContactName {
		t.Errorf("expected contact name '%s', got '%s'", req.ContactName, response.ContactName)
	}

	if response.Status != "Pending" {
		t.Errorf("expected status 'Pending', got '%s'", response.Status)
	}

	// 测试重复创建（应该失败）
	_, err = service.CreateFranchiseIntention(testMember.ID, req)
	if err == nil {
		t.Error("expected error for duplicate franchise intention")
	}

	// 测试空联系人姓名
	emptyReq := contracts.CreateFranchiseIntentionRequest{
		ContactName:      "",
		ContactPhone:     "13800138000",
		IntendedLocation: "北京市朝阳区",
	}

	_, err = service.CreateFranchiseIntention("another-member-id", emptyReq)
	if err == nil {
		t.Error("expected error for empty contact name")
	}
}

func TestMemberService_GetMemberInfo(t *testing.T) {
	db := setupServiceTestDB(t)
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)
	service := NewMemberService(memberRepo, franchiseRepo)

	// 创建测试数据
	testMember := createServiceTestMember(t, db)

	// 创建加盟意向
	intention := &models.FranchiseIntention{
		ID:               "service-test-intention-1",
		MemberID:         testMember.ID,
		ContactName:      "李四",
		ContactPhone:     "13900139000",
		IntendedLocation: "上海市",
		Status:           "Pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := franchiseRepo.Create(intention)
	if err != nil {
		t.Fatalf("failed to create test intention: %v", err)
	}

	// 测试获取会员信息
	response, err := service.GetMemberInfo(testMember.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.ID != testMember.ID {
		t.Errorf("expected ID '%s', got '%s'", testMember.ID, response.ID)
	}

	if response.Nickname != testMember.Nickname {
		t.Errorf("expected nickname '%s', got '%s'", testMember.Nickname, response.Nickname)
	}

	if len(response.FranchiseIntentions) != 1 {
		t.Errorf("expected 1 franchise intention, got %d", len(response.FranchiseIntentions))
	}

	if response.FranchiseIntentions[0].ContactName != intention.ContactName {
		t.Errorf("expected contact name '%s', got '%s'", intention.ContactName, response.FranchiseIntentions[0].ContactName)
	}

	// 测试不存在的用户
	_, err = service.GetMemberInfo("non-existent-id")
	if err == nil {
		t.Error("expected error for non-existent member")
	}
}

func TestMemberService_ValidateMemberExists(t *testing.T) {
	db := setupServiceTestDB(t)
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)
	service := NewMemberService(memberRepo, franchiseRepo)

	// 创建测试数据
	testMember := createServiceTestMember(t, db)

	// 测试存在的用户
	err := service.ValidateMemberExists(testMember.ID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 测试不存在的用户
	err = service.ValidateMemberExists("non-existent-id")
	if err == nil {
		t.Error("expected error for non-existent member")
	}
}

func TestMemberService_UpdateFranchiseIntentionStatus(t *testing.T) {
	db := setupServiceTestDB(t)
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)
	service := NewMemberService(memberRepo, franchiseRepo)

	// 创建测试数据
	testMember := createServiceTestMember(t, db)
	intention := &models.FranchiseIntention{
		ID:               "test-status-intention-1",
		MemberID:         testMember.ID,
		ContactName:      "王五",
		ContactPhone:     "13700137000",
		IntendedLocation: "深圳市",
		Status:           "Pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := franchiseRepo.Create(intention)
	if err != nil {
		t.Fatalf("failed to create test intention: %v", err)
	}

	// 测试有效状态更新
	err = service.UpdateFranchiseIntentionStatus(intention.ID, "Approved")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 测试无效状态
	err = service.UpdateFranchiseIntentionStatus(intention.ID, "InvalidStatus")
	if err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestMemberService_generateFranchiseIntentionID(t *testing.T) {
	db := setupServiceTestDB(t)
	memberRepo := repositories.NewMemberRepository(db)
	franchiseRepo := repositories.NewFranchiseIntentionRepository(db)
	service := NewMemberService(memberRepo, franchiseRepo)

	id := service.generateFranchiseIntentionID()

	if id == "" {
		t.Error("expected non-empty ID")
	}

	if len(id) != 11 { // "fi-" + 8 characters
		t.Errorf("expected ID length 11, got %d", len(id))
	}

	if id[:3] != "fi-" {
		t.Errorf("expected ID to start with 'fi-', got '%s'", id[:3])
	}

	// 测试生成的ID是唯一的
	id2 := service.generateFranchiseIntentionID()
	if id == id2 {
		t.Error("expected different IDs to be generated")
	}
}
