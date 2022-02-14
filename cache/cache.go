// The cache package provides content caching and a mechanism to prevent cache stampedes.
//
// Cache stampedes occur when two or more requests for expensive operations are received
// at the same time, before the content is cached.
// The requests find the content is not in the cache, and then they all do the expensive
// operation.
//
// With this package, cache stampedes are avoided by requiring handlers to allocate and lock
// a cache key before they are permitted to access the cache.
//
// General usage:
// 1. application calls New to set up the cache manager
// 2. on each request:
//	a. generate a cache key based on the request
//	b. allocate an Entry for this cache key
//	c. lock the Entry for the duration of cache operations
//	d. unlock the Entry after cache operations are complete, but before writing to client
//	3. release the cache key entry
//
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
)

// A Manager manages the underlying cache and the dynamic set of Entries.
type Manager struct {
	cache      *cache.Cache      // underlying cache
	sync.Mutex                   // protexts operations on locks and references below
	entries    map[string]*Entry // cache access manager for each key; map index is the key
	references map[string]int    // reference counts for each key; map index is the key
}

// An Entry manages cache access and locking for a single cache key.
type Entry struct {
	key        string   // the cache key
	manager    *Manager // our "parent" manager
	sync.Mutex          // mutex to serialise cache operations related to this key
}

// New sets up a new cache and lock manager.
func New(ttl time.Duration, megabytes int) (*Manager, error) {
	// configure bigcache
	config := bigcache.DefaultConfig(ttl)
	config.HardMaxCacheSize = megabytes

	// create bigcache client
	bigcacheClient, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	// use bigcache client as a gocache store
	bigcacheStore := store.NewBigcache(bigcacheClient, nil)

	// create single stage cache using bigcache
	cacheManager := cache.New(bigcacheStore)

	return &Manager{
		cache:      cacheManager,
		entries:    map[string]*Entry{},
		references: map[string]int{},
	}, nil
}

// AllocateEntry returns an object that may be locked to serialise
// access to cache entries for the given key.
// The returned object is also used to Get and Set cache values
// for the key.
// A reference count is incremented on each call, so you must call
// Free when you are done with this object.
func (cm *Manager) AllocateEntry(key string) *Entry {
	cm.Lock()
	defer cm.Unlock()

	entry, ok := cm.entries[key]
	if !ok {
		entry = &Entry{
			key:     key,
			manager: cm,
		}
		cm.entries[key] = entry
		cm.references[key] = 0
	}
	cm.references[key]++
	return entry
}

// Free gives a key serialiser back and decrements its reference count.
// If the reference count drops to zero, the entire serialiser is freed.
// Do not try to use the serialiser once it has been freed.
// Do not call Free if you still have a lock on this serialiser.
func (entry *Entry) Free() {
	key := entry.key
	cm := entry.manager // our parent CacheManager

	cm.Lock()
	defer cm.Unlock()

	cm.references[key]--
	if cm.references[key] == 0 {
		delete(cm.entries, key)
		delete(cm.references, key)
	}
}

// Get retrieves a value from the cache for key.
func (entry *Entry) Get(ctx context.Context) ([]byte, error) {
	v, err := entry.manager.cache.Get(ctx, entry.key)
	if err != nil {
		return nil, err
	}
	return v.([]byte), nil
}

// Set saves a new value in the cache for key.
func (entry *Entry) Set(ctx context.Context, value []byte) error {
	return entry.manager.cache.Set(ctx, entry.key, value, nil)
}
