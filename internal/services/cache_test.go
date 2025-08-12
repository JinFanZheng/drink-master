package services

import (
	"testing"
	"time"
)

func TestNewCacheManager(t *testing.T) {
	cm := NewCacheManager()

	if cm == nil {
		t.Fatal("Cache manager should not be nil")
	}

	if cm.data == nil {
		t.Error("Cache data map should not be nil")
	}

	if cm.ttl != 24*time.Hour {
		t.Errorf("Expected TTL 24h, got %v", cm.ttl)
	}
}

func TestCacheManager_SetAndGetLoginStatus(t *testing.T) {
	cm := NewCacheManager()
	memberID := "test_member_id"
	token := "test_token"

	// Set login status
	cm.SetLoginStatus(memberID, token)

	// Get login status
	retrievedToken, exists := cm.GetLoginStatus(memberID)

	if !exists {
		t.Error("Login status should exist")
	}

	if retrievedToken != token {
		t.Errorf("Expected token %s, got %s", token, retrievedToken)
	}
}

func TestCacheManager_GetNonExistentLoginStatus(t *testing.T) {
	cm := NewCacheManager()

	token, exists := cm.GetLoginStatus("non_existent_member")

	if exists {
		t.Error("Login status should not exist")
	}

	if token != "" {
		t.Error("Token should be empty for non-existent member")
	}
}

func TestCacheManager_RemoveLoginStatus(t *testing.T) {
	cm := NewCacheManager()
	memberID := "test_member_id"
	token := "test_token"

	// Set login status
	cm.SetLoginStatus(memberID, token)

	// Verify it exists
	_, exists := cm.GetLoginStatus(memberID)
	if !exists {
		t.Error("Login status should exist before removal")
	}

	// Remove login status
	cm.RemoveLoginStatus(memberID)

	// Verify it's removed
	_, exists = cm.GetLoginStatus(memberID)
	if exists {
		t.Error("Login status should not exist after removal")
	}
}

func TestCacheManager_CleanupExpired(t *testing.T) {
	cm := NewCacheManager()

	// Set a login status
	memberID := "test_member_cleanup"
	token := "test_token"
	cm.SetLoginStatus(memberID, token)

	// Manually set an expired entry
	expiredMemberID := "expired_member"
	expiredToken := "expired_token"

	// Add entry and manually mark it as expired
	cm.mutex.Lock()
	cm.data["login:"+expiredMemberID] = CacheEntry{
		Value:     expiredToken,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Already expired
	}
	cm.mutex.Unlock()

	// Verify both entries exist initially
	if _, exists := cm.GetLoginStatus(memberID); !exists {
		t.Error("Active login status should exist")
	}

	// Manually trigger cleanup
	cm.cleanupExpired()

	// Verify active entry still exists
	if _, exists := cm.GetLoginStatus(memberID); !exists {
		t.Error("Active login status should still exist after cleanup")
	}

	// The expired entry should be removed but since GetLoginStatus checks expiry anyway,
	// we need to check the underlying data map
	cm.mutex.RLock()
	_, expiredExists := cm.data["login:"+expiredMemberID]
	cm.mutex.RUnlock()

	if expiredExists {
		t.Error("Expired login status should be removed after cleanup")
	}
}

func TestCacheManager_ExpiredEntry(t *testing.T) {
	cm := NewCacheManager()
	cm.ttl = 1 * time.Millisecond // Very short TTL for testing

	memberID := "test_member_id"
	token := "test_token"

	// Set login status
	cm.SetLoginStatus(memberID, token)

	// Wait for expiration
	time.Sleep(2 * time.Millisecond)

	// Try to get expired entry
	_, exists := cm.GetLoginStatus(memberID)

	if exists {
		t.Error("Expired login status should not exist")
	}
}
