package cache

import "net/http"

// cacheKey builds a cache key from an incoming HTTP request struct.
// Right now it only looks at RequestURI, but it should probably
// also look at certain headers to do with language, content-type,
// encoding, etc.
func CacheKey(req *http.Request) string {
	return req.RequestURI
}
