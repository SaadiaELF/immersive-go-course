package main

import "fmt"

type Cache struct {
	storage  map[int]any
	capacity int
}

func NewCache(capacity int) *Cache {
	return &Cache{
		storage:  make(map[int]any),
		capacity: capacity,
	}
}

func main() {
	cache := NewCache(10)
	fmt.Println(cache)
}
