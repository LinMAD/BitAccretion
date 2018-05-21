package cache

import (
	"sync"
)

// MemoryCache temporary storage to keep some objects in memory
type MemoryCache struct {
	// bucket a collection of stored objects
	bucket map[string]item
	sync.RWMutex
}

// item is a simple structure of a subject which will be stored in memory
type item struct {
	// Placed object in memory
	object interface{}
}

// Boot prepares memory caching
func Boot() *MemoryCache {
	return &MemoryCache{
		bucket: make(map[string]item),
	}
}

// Add data under given key
// If you use same key it(data) will be rewritten in memory
func (mc *MemoryCache) Add(key string, data interface{}) {
	mc.Lock()
	if _, found := mc.bucket[key]; found {
		delete(mc.bucket, key)
	}
	mc.bucket[key] = item{object: data}
	mc.Unlock()
}

// Get data by key if it's expired then will be returned nil
// Every time when data gets from memory, then TTL will be updated
func (mc *MemoryCache) Get(key string) interface{} {
	mc.RLock()

	// Try get cached object
	item, inStack := mc.bucket[key]
	if !inStack {
		mc.RUnlock()

		return nil
	}
	mc.RUnlock()

	return item.object
}

// DeleteBy key object in cache
func (mc *MemoryCache) DeleteBy(key string) (isDeleted bool) {
	mc.Lock()
	isDeleted = mc.delete(key)
	mc.Unlock()
	return
}

// delete data from bucket
func (mc *MemoryCache) delete(key string) bool {
	if _, found := mc.bucket[key]; found {
		delete(mc.bucket, key)

		return true
	}

	return false
}
