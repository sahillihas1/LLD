package main

import (
	"container/list"
	"fmt"
)

type EvictionStrategy interface {
	Evict(cache *LRUCache)
}

// LRUEviction implements the LRU eviction strategy
type LRUEviction struct{}

func (l *LRUEviction) Evict(cache *LRUCache) {
	if el := cache.evictionList.Back(); el != nil {
		cache.evictionList.Remove(el)
		delete(cache.data, el.Value.(*Entry).key)
	}
}

// Cache interface defines basic cache operations
type Cache interface {
	Put(key string, value interface{})
	Get(key string) (interface{}, bool)
	SetEvictionStrategy(strategy EvictionStrategy)
}

// LRUCache implements the LRU Cache
type LRUCache struct {
	capacity         int
	data             map[string]*list.Element
	evictionList     *list.List
	evictionStrategy EvictionStrategy
}

// Entry represents a key-value pair in the cache
type Entry struct {
	key   string
	value interface{}
}

// Constructor for LRUCache
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity:     capacity,
		data:         make(map[string]*list.Element),
		evictionList: list.New(),
	}
}

// Put adds an item to the cache
func (c *LRUCache) Put(key string, value interface{}) {
	if el, ok := c.data[key]; ok {
		c.evictionList.MoveToFront(el)
		el.Value.(*Entry).value = value
		return
	}
	if len(c.data) >= c.capacity {
		c.evictionStrategy.Evict(c)
	}
	el := c.evictionList.PushFront(&Entry{key, value})
	c.data[key] = el
}

// Get retrieves an item from the cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	if el, ok := c.data[key]; ok {
		c.evictionList.MoveToFront(el)
		return el.Value.(*Entry).value, true
	}
	return nil, false
}

// SetEvictionStrategy sets the eviction strategy for the cache
func (c *LRUCache) SetEvictionStrategy(strategy EvictionStrategy) {
	c.evictionStrategy = strategy
}

// ============================ Factory Pattern (Cache Factory) ============================

// CacheFactory creates caches based on type
type CacheFactory struct{}

func (cf *CacheFactory) CreateCache(cacheType string, capacity int) Cache {
	switch cacheType {
	case "LRU":
		return NewLRUCache(capacity)
	default:
		return nil
	}
}

// ============================ Builder Pattern (Cache Builder) ============================

// CacheBuilder helps construct a cache
type CacheBuilder struct {
	cache Cache
}

func NewCacheBuilder() *CacheBuilder {
	return &CacheBuilder{}
}

func (cb *CacheBuilder) SetCacheType(cacheType string, capacity int) *CacheBuilder {
	cb.cache = NewLRUCache(capacity)
	return cb
}

func (cb *CacheBuilder) Build() Cache {
	return cb.cache
}

// ============================ Main Function (Testing) ============================

func main() {
	// Using Factory Pattern to create a cache
	cacheFactory := &CacheFactory{}
	cache := cacheFactory.CreateCache("LRU", 3)

	// Set eviction strategy (Strategy Pattern)
	evictionStrategy := &LRUEviction{}
	cache.SetEvictionStrategy(evictionStrategy)

	// Using Builder Pattern to construct cache
	cacheBuilder := NewCacheBuilder()
	cache = cacheBuilder.SetCacheType("LRU", 3).Build()

	// Adding elements
	cache.Put("A", 1)
	cache.Put("B", 2)
	cache.Put("C", 3)
	fmt.Println(cache.Get("A")) // Output: 1, true

	// Exceeding capacity to trigger eviction
	cache.Put("D", 4)
	fmt.Println(cache.Get("B")) // Output: nil, false (Evicted)
}
