package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFranchiseIntention_TableName(t *testing.T) {
	fi := FranchiseIntention{}
	expected := "franchise_intentions"
	assert.Equal(t, expected, fi.TableName())
}

func TestFranchiseIntention_BeforeCreate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	memberID := "member-123"
	company := "测试公司"
	name := "张三"
	mobile := "13800138000"
	area := "北京市"
	remark := "测试备注"

	fi := &FranchiseIntention{
		ID:        "test-id",
		MemberId:  &memberID,
		Company:   &company,
		Name:      &name,
		Mobile:    &mobile,
		Area:      &area,
		Remark:    &remark,
		IsHandled: BitBool(0),
		Version:   1,
		CreatedOn: time.Now(),
	}

	err = fi.BeforeCreate(db)
	assert.NoError(t, err)
}

func TestFranchiseIntention_BeforeUpdate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	memberID := "member-123"
	now := time.Now()

	fi := &FranchiseIntention{
		ID:        "test-id",
		MemberId:  &memberID,
		IsHandled: BitBool(1),
		Version:   1,
		CreatedOn: now,
		UpdatedOn: &now,
	}

	err = fi.BeforeUpdate(db)
	assert.NoError(t, err)
}

func TestFranchiseIntention_GetHandledStatus(t *testing.T) {
	tests := []struct {
		name      string
		isHandled BitBool
		expected  bool
	}{
		{
			name:      "handled status true",
			isHandled: BitBool(1),
			expected:  true,
		},
		{
			name:      "handled status false",
			isHandled: BitBool(0),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fi := &FranchiseIntention{
				IsHandled: tt.isHandled,
			}
			assert.Equal(t, tt.expected, fi.GetHandledStatus())
		})
	}
}

func TestFranchiseIntention_IsPending(t *testing.T) {
	tests := []struct {
		name      string
		isHandled BitBool
		expected  bool
	}{
		{
			name:      "is pending - not handled",
			isHandled: BitBool(0),
			expected:  true,
		},
		{
			name:      "is not pending - handled",
			isHandled: BitBool(1),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fi := &FranchiseIntention{
				IsHandled: tt.isHandled,
			}
			assert.Equal(t, tt.expected, fi.IsPending())
		})
	}
}

func TestFranchiseIntention_ModelStructure(t *testing.T) {
	memberID := "member-123"
	company := "测试公司"
	name := "张三"
	mobile := "13800138000"
	area := "北京市"
	remark := "测试备注"
	now := time.Now()

	fi := FranchiseIntention{
		ID:        "fi-123",
		MemberId:  &memberID,
		Company:   &company,
		Name:      &name,
		Mobile:    &mobile,
		Area:      &area,
		Remark:    &remark,
		IsHandled: BitBool(0),
		Version:   1,
		CreatedOn: now,
		UpdatedOn: &now,
	}

	assert.Equal(t, "fi-123", fi.ID)
	assert.NotNil(t, fi.MemberId)
	assert.Equal(t, "member-123", *fi.MemberId)
	assert.NotNil(t, fi.Company)
	assert.Equal(t, "测试公司", *fi.Company)
	assert.NotNil(t, fi.Name)
	assert.Equal(t, "张三", *fi.Name)
	assert.NotNil(t, fi.Mobile)
	assert.Equal(t, "13800138000", *fi.Mobile)
	assert.NotNil(t, fi.Area)
	assert.Equal(t, "北京市", *fi.Area)
	assert.NotNil(t, fi.Remark)
	assert.Equal(t, "测试备注", *fi.Remark)
	assert.Equal(t, BitBool(0), fi.IsHandled)
	assert.Equal(t, int64(1), fi.Version)
	assert.Equal(t, now, fi.CreatedOn)
	assert.NotNil(t, fi.UpdatedOn)
}
