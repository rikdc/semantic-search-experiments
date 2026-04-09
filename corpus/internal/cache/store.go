package cache

import (
	"sync"
	"time"
)

// entry holds a cached value alongside its expiry time.
type entry struct {
	value     any
	expiresAt time.Time
}

// Store is a thread-safe in-memory key-value cache with optional TTL per entry.
type Store struct {
	mu   sync.RWMutex
	data map[string]entry
}

// NewStore creates an empty Store ready for use.
func NewStore() *Store {
	return &Store{data: make(map[string]entry)}
}

// Set stores a value under the given key. If ttl is zero the entry never expires.
func Set(s *Store, key string, value any, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	exp := time.Time{}
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	s.data[key] = entry{value: value, expiresAt: exp}
}

// Get retrieves the value for key. Returns (value, true) on a hit and (nil, false)
// on a miss or if the entry has expired.
func Get(s *Store, key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.data[key]
	if !ok {
		return nil, false
	}
	if !e.expiresAt.IsZero() && time.Now().After(e.expiresAt) {
		return nil, false
	}
	return e.value, true
}

// Delete removes the entry for key. It is a no-op if the key does not exist.
func Delete(s *Store, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Flush removes all entries from the store.
func Flush(s *Store) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]entry)
}
