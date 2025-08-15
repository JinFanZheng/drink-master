package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/models"
)

func setupFranchiseTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = models.AutoMigrate(db)
	assert.NoError(t, err)

	return db
}

func createTestFranchiseIntention(t *testing.T, db *gorm.DB, memberID string) *models.FranchiseIntention {
	company := "测试公司"
	name := "张三"
	mobile := "13800138000"
	area := "北京市朝阳区"
	remark := "测试备注"
	now := time.Now()

	intention := &models.FranchiseIntention{
		ID:        "test-franchise-1",
		MemberId:  &memberID,
		Company:   &company,
		Name:      &name,
		Mobile:    &mobile,
		Area:      &area,
		Remark:    &remark,
		IsHandled: models.BitBool(0),
		Version:   1,
		CreatedOn: now,
		UpdatedOn: &now,
	}

	err := db.Create(intention).Error
	assert.NoError(t, err)

	return intention
}

func TestFranchiseIntentionRepository_Create(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	memberID := "member-123"
	company := "新公司"
	name := "李四"
	mobile := "13900139000"
	area := "上海市"
	remark := "创建测试"
	now := time.Now()

	intention := &models.FranchiseIntention{
		ID:        "create-test-1",
		MemberId:  &memberID,
		Company:   &company,
		Name:      &name,
		Mobile:    &mobile,
		Area:      &area,
		Remark:    &remark,
		IsHandled: models.BitBool(0),
		Version:   1,
		CreatedOn: now,
	}

	err := repo.Create(intention)
	assert.NoError(t, err)

	// 验证创建结果
	var result models.FranchiseIntention
	err = db.First(&result, "id = ?", intention.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, intention.ID, result.ID)
	assert.Equal(t, memberID, *result.MemberId)
}

func TestFranchiseIntentionRepository_GetByMemberID(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	memberID := "member-456"
	_ = createTestFranchiseIntention(t, db, memberID)

	intentions, err := repo.GetByMemberID(memberID)
	assert.NoError(t, err)
	assert.Len(t, intentions, 1)
	assert.Equal(t, memberID, *intentions[0].MemberId)
}

func TestFranchiseIntentionRepository_UpdateDirectly(t *testing.T) {
	db := setupFranchiseTestDB(t)
	repo := NewFranchiseIntentionRepository(db)

	memberID := "member-789"
	intention := createTestFranchiseIntention(t, db, memberID)

	// 直接更新处理状态
	intention.IsHandled = models.BitBool(1)
	err := repo.Update(intention)
	assert.NoError(t, err)

	// 验证更新结果
	var result models.FranchiseIntention
	err = db.First(&result, "id = ?", intention.ID).Error
	assert.NoError(t, err)
	assert.True(t, result.GetHandledStatus())
}

func TestFranchiseIntentionRepository_GetPendingDirectly(t *testing.T) {
	db := setupFranchiseTestDB(t)
	_ = NewFranchiseIntentionRepository(db)

	// 创建待处理的意向
	memberID1 := "member-pending-1"
	_ = createTestFranchiseIntention(t, db, memberID1)

	// 创建已处理的意向
	memberID2 := "member-handled-2"
	handledIntention := createTestFranchiseIntention(t, db, memberID2)
	handledIntention.IsHandled = models.BitBool(1)
	db.Save(handledIntention)

	// 查询待处理的意向
	var pendingIntentions []models.FranchiseIntention
	err := db.Where("is_handled = ?", 0).Find(&pendingIntentions).Error
	assert.NoError(t, err)
	assert.Len(t, pendingIntentions, 1)
	assert.True(t, pendingIntentions[0].IsPending())
}
