package handlers

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewBaseHandler(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	handler := NewBaseHandler(db)
	
	if handler.db != db {
		t.Error("Expected handler.db to be set correctly")
	}
}