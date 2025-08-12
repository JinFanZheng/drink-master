package services

import (
	"fmt"
	"sync"
	"time"
)

// CacheEntry represents a cache entry with expiration
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// CacheManager handles login status caching (in-memory for now)
// In production, this should use Redis or similar distributed cache
type CacheManager struct {
	data  map[string]CacheEntry
	mutex sync.RWMutex
	ttl   time.Duration
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	cm := &CacheManager{
		data: make(map[string]CacheEntry),
		ttl:  24 * time.Hour, // Default TTL of 24 hours
	}

	// Start cleanup goroutine
	go cm.cleanup()

	return cm
}

// SetLoginStatus caches login status for a member
func (c *CacheManager) SetLoginStatus(memberID, token string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[fmt.Sprintf("login:%s", memberID)] = CacheEntry{
		Value:     token,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// GetLoginStatus retrieves login status for a member
func (c *CacheManager) GetLoginStatus(memberID string) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.data[fmt.Sprintf("login:%s", memberID)]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return "", false
	}

	token, ok := entry.Value.(string)
	return token, ok
}

// RemoveLoginStatus removes login status for a member
func (c *CacheManager) RemoveLoginStatus(memberID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, fmt.Sprintf("login:%s", memberID))
}

// cleanup runs a background process to clean up expired entries
func (c *CacheManager) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpired()
	}
}

// cleanupExpired removes expired entries from cache
func (c *CacheManager) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.After(entry.ExpiresAt) {
			delete(c.data, key)
		}
	}
}
