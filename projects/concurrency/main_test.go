package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPut(t *testing.T) {
	cache := NewCache(2)
	t.Run("Case1/Below capacity", func(t *testing.T) {
		cache.Put(1, "one")
		cache.Put(2, "two")
		_, exists := cache.storage[1]
		require.True(t, exists)
		_, exists = cache.storage[2]
		require.True(t, exists)
	})
	t.Run("Case1/Beyond capacity", func(t *testing.T) {
		cache.Put(1, "one")
		cache.Put(2, "two")
		cache.Put(3, "three")
		_, exists := cache.storage[1]
		require.False(t, exists)
		_, exists = cache.storage[2]
		require.True(t, exists)
		_, exists = cache.storage[3]
		require.True(t, exists)
	})
}

func TestGet(t *testing.T) {
	cache := NewCache(2)
	cache.Put(1, "one")
	cache.Put(2, "two")
	t.Run("Case1/Exists", func(t *testing.T) {
		value, exists := cache.Get(1)
		require.True(t, exists)
		require.Equal(t, "one", value)
		// check the most recently used key is 1
		require.True(t, cache.storage[1].lastUsed.After(cache.storage[2].lastUsed))
	})
	t.Run("Case2/Does not exist", func(t *testing.T) {
		value, exists := cache.Get(3)
		require.False(t, exists)
		require.Nil(t, value)
	})
}

// test which tries to concurrently use the map from a bunch of goroutines
func TestConcurrency(t *testing.T) {
	value, exist := any(nil), false
	cache := NewCache(2)
	done := make(chan bool)

	// Failing test as the Get is called before the Put is done , because the cache is not thread safe
	go func() {
		cache.Put(1, "one")
		fmt.Println("Put 1")
		done <- true
	}()
	go func() {
		value, exist = cache.Get(1)
		fmt.Println("Get 1")
		done <- true
	}()
	for i := 0; i < 2; i++ {
		<-done
	}
	require.True(t, exist)
	require.Equal(t, "one", value)

}
