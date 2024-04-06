package main

import (
	"fmt"
	"time"
)

type Cache struct {
	capacity int
	storage  map[int]entry
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
	fmt.Println(cache) 
	cache.Put(3, "three")
	fmt.Println(cache)
	cache.Get(2)
	cache.Put(4, "four")

	fmt.Println(cache)
}
