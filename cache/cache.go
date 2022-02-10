package cache

import (
	"context"
	"net/http"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
)

type Cache struct {
	cache *cache.Cache
}

func New(ttl time.Duration, megabytes int) (*Cache, error) {

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

	return &Cache{
		cache: cacheManager,
	}, nil
}

func (c *Cache) Get(ctx context.Context, req *http.Request) ([]byte, error) {
	v, err := c.cache.Get(ctx, cacheKey(req))
	if err != nil {
		return nil, err
	}
	return v.([]byte), nil
}

func (c *Cache) Set(ctx context.Context, req *http.Request, body []byte) error {
	return c.cache.Set(ctx, cacheKey(req), body, nil)
}

func cacheKey(req *http.Request) string {
	return req.RequestURI
}
