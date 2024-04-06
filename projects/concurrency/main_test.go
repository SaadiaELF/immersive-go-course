package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func (c *Cache) GetKeys() []int {
	keys := make([]int, 0, len(c.storage))
	for key := range c.storage {
		keys = append(keys, key)
	}
	return keys
}
func TestPut(t *testing.T) {
	cache := NewCache(2)
	t.Run("Case1/Below capacity", func(t *testing.T) {
		cache.Put(1, "one")
		cache.Put(2, "two")
		actual := cache.GetKeys()
		require.Equal(t, []int{1, 2}, actual)
	})
	t.Run("Case1/Beyond capacity", func(t *testing.T) {
		cache.Put(1, "one")
		cache.Put(2, "two")
		cache.Put(3, "three")
		actual := cache.GetKeys()
		require.Equal(t, []int{2, 3}, actual)
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
