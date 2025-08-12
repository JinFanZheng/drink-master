package models

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFranchiseIntention_TableName(t *testing.T) {
	fi := FranchiseIntention{}
	expected := "franchise_intentions"

	if fi.TableName() != expected {
		t.Errorf("Expected table name %s, got %s", expected, fi.TableName())
	}
}

func TestFranchiseIntention_BeforeCreate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 自动迁移
	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	fi := &FranchiseIntention{
		ID:               "test-id",
		MemberID:         "member-123",
		ContactName:      "张三",
		ContactPhone:     "13800138000",
		IntendedLocation: "北京市",
		Status:           "Pending",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 测试BeforeCreate钩子
	err = fi.BeforeCreate(db)
	if err != nil {
		t.Errorf("BeforeCreate should not return error: %v", err)
	}
}

func TestFranchiseIntention_BeforeUpdate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	fi := &FranchiseIntention{
		Status: "InvalidStatus",
	}

	// 测试BeforeUpdate钩子 - 无效状态会被重置为Pending
	err = fi.BeforeUpdate(db)
	if err != nil {
		t.Errorf("BeforeUpdate should not return error: %v", err)
	}

	if fi.Status != "Pending" {
		t.Errorf("Expected status to be reset to 'Pending', got %s", fi.Status)
	}

	// 测试有效状态
	fi.Status = "Approved"
	err = fi.BeforeUpdate(db)
	if err != nil {
		t.Errorf("BeforeUpdate should not return error: %v", err)
	}

	if fi.Status != "Approved" {
		t.Errorf("Expected status to remain 'Approved', got %s", fi.Status)
	}
}

func TestFranchiseIntention_StatusMethods(t *testing.T) {
	testCases := []struct {
		status     string
		isPending  bool
		isApproved bool
		isRejected bool
	}{
		{"Pending", true, false, false},
		{"Approved", false, true, false},
		{"Rejected", false, false, true},
		{"Unknown", false, false, false},
	}

	for _, tc := range testCases {
		fi := FranchiseIntention{Status: tc.status}

		if fi.IsPending() != tc.isPending {
			t.Errorf("Status %s: expected IsPending %v, got %v", tc.status, tc.isPending, fi.IsPending())
		}

		if fi.IsApproved() != tc.isApproved {
			t.Errorf("Status %s: expected IsApproved %v, got %v", tc.status, tc.isApproved, fi.IsApproved())
		}

		if fi.IsRejected() != tc.isRejected {
			t.Errorf("Status %s: expected IsRejected %v, got %v", tc.status, tc.isRejected, fi.IsRejected())
		}
	}
}
