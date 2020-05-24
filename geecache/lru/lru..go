package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int64
	nBytes int64

	ll *list.List

	cache map[string]*list.Element
	OnEnvicted func(key string, value Value)
}

type entry struct {
	key string
	value Value
}

type Value interface {
	Len() int
}

func NewCache(maxBytes int64, onEnvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		OnEnvicted:  onEnvicted,
		cache: make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, ok
	}

	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()

	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEnvicted != nil {
			c.OnEnvicted(kv.key, kv.value)
		}
	}

	return
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
 }

 func (c *Cache) Len() int {
 	return c.ll.Len()
 }