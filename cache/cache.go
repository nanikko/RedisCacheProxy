package cache

import (
	"github.com/hashicorp/golang-lru/simplelru"
    "log"
    "time"
)

type Cache interface {
	Get(string) (bool, string)
	Add(string, string)
}

type CacheValue struct {
	v string
	lastUpdatedTime int64
}

type MinimalCache struct {
	lru simplelru.LRUCache
	expiry int64
}

func NewMinimalCache(cacheSize int, expiredTime int64) (*MinimalCache, error) {
	l, err := simplelru.NewLRU(cacheSize, nil)
	c := &MinimalCache{
		lru : l,
		expiry : expiredTime,
	}
	return c, err
}

func (c *MinimalCache) Get(key string) (bool, string) {
	// Lru get() will move the entry to front
	if val,ok := c.lru.Get(key); ok {
		// check if expired
		cachedv := val.(CacheValue)

		if (time.Now().Unix() - cachedv.lastUpdatedTime) >= c.expiry {
			// expired
			log.Printf("Expired Key:%s in cache", key)
			c.lru.Remove(key)
			return false, ""
		}

		log.Println("Non-expired Key:", key)
		return ok, cachedv.v
	}

	return false, ""
}

func (c *MinimalCache) Add(k, v string) {
	log.Printf("Added: Value:%s against Key:%s in cache", v, k)
	t := time.Now().Unix()
	cVal := CacheValue{v, t}
	c.lru.Add(k, cVal)
}



