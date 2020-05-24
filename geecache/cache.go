package geecache

import (
	"gee/geecache/lru"
	"sync"
)

type cache struct {
	m sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

func (c *cache) Add(key string, value ByteView) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}

	c.lru.Add(key, value)
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.lru == nil {
		return
	}

	if val, ok := c.lru.Get(key); ok {
		return val.(ByteView), ok
	}

	return
}


