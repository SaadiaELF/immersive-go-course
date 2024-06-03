package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	capacity int
	storage  map[int]entry
	mu       sync.Mutex
}

type entry struct {
	value    any
	lastUsed time.Time
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		storage:  make(map[int]entry, capacity),
	}
}

// Put adds the value to the cache, and returns a boolean to indicate whether a value already existed in the cache for that key.
// If there was previously a value, it replaces that value with this one.
// Any Put counts as a refresh in terms of LRU tracking.
func (c *Cache) Put(key int, value any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Check if the entry already exists in the cache
	_, exists := c.storage[key]
	// Update the entry with the new value and current time
	c.storage[key] = entry{value: value, lastUsed: time.Now()}
	// If the cache is at capacity, remove the least recently used (LRU) entry
	if len(c.storage) > c.capacity {
		var lruKey int
		var lruTime time.Time
		for k, e := range c.storage {
			if lruTime.IsZero() || e.lastUsed.Before(lruTime) {
				lruKey = k
				lruTime = e.lastUsed
			}
		}

		delete(c.storage, lruKey)
	}

	// Return true if the entry already existed, false otherwise
	return exists
}

// Get returns the value associated with the passed key, and a boolean to indicate whether a value was known or not. If not, nil is returned as the value.
// Any Get counts as a refresh in terms of LRU tracking.
func (c *Cache) Get(key int) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Check if the entry exists in the cache
	entry, exists := c.storage[key]
	// If the entry exists, update the last used time
	if exists {
		entry.lastUsed = time.Now()
		c.storage[key] = entry
	}
	// Return the value and whether it existed or not
	return entry.value, exists
}

func main() {
	cache := NewCache(2)
	cache.Put(1, "one")
	cache.Put(2, "two")
	fmt.Println(cache) // Print cache with 1 and 2
	cache.Put(3, "three")
	fmt.Println(cache) // Print cache with 2 and 3 as 1 is removed due to LRU
	cache.Get(2)
	fmt.Println(cache) // Print cache with 3 and 2 as 2 is refreshed
	cache.Put(4, "four")
	fmt.Println(cache) // Print cache with 2 and 4 as 3 is removed due to LRU
}
