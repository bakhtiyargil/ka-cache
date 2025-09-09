package cache

import (
	"errors"
	"ka-cache/logger"
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Put(key K, value V, ttl int64) error
	Get(key K) (*Entry[K, V], bool)
	TTL(key K) time.Duration
}

type Entry[K comparable, V any] struct {
	key       K
	Value     V
	expiresAt time.Time
	next      *Entry[K, V]
	prev      *Entry[K, V]
}

type LruCache[K comparable, V any] struct {
	cacheMap    map[K]*Entry[K, V]
	logger      logger.Logger
	rwMutex     sync.RWMutex
	capacity    int
	cleanupStop chan bool
	head        *Entry[K, V]
	tail        *Entry[K, V]
}

func NewLruCache[K comparable, V any](cap int, logger logger.Logger) SelfClearingCache[K, V] {
	newCacheMap := make(map[K]*Entry[K, V], cap)
	cache := LruCache[K, V]{
		cacheMap: newCacheMap,
		capacity: cap,
		logger:   logger,
		head:     nil,
		tail:     nil,
	}
	return &cache
}

func (c *LruCache[K, V]) Put(key K, value V, ttl int64) (err error) {
	strKey, key := c.checkStrKey(key)
	if len(strKey) == 0 {
		return c.putAny(key, value, ttl)
	}
	return c.putAny(any(strKey).(K), value, ttl)
}

func (c *LruCache[K, V]) putAny(key K, value V, ttl int64) error {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	existingNode, ok := c.cacheMap[key]
	if ok {
		existingNode.Value = value
		if err := c.setExpirationTime(existingNode, ttl); err != nil {
			return err
		}
		c.unlink(existingNode)
		c.linkFirst(existingNode)
	} else {
		if len(c.cacheMap) >= c.capacity {
			c.deleteAndUnlink(c.tail)
		}
		var newEntry Entry[K, V]
		newEntry.key = key
		newEntry.Value = value
		if err := c.setExpirationTime(&newEntry, ttl); err != nil {
			return err
		}
		c.cacheMap[key] = &newEntry
		c.linkFirst(&newEntry)
	}
	return nil
}

func (c *LruCache[K, V]) Get(key K) (*Entry[K, V], bool) {
	strKey, key := c.checkStrKey(key)
	if len(strKey) == 0 {
		return c.getAny(key)
	}
	return c.getAny(any(strKey).(K))
}

func (c *LruCache[K, V]) getAny(key K) (*Entry[K, V], bool) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	cacheEntry, ok := c.cacheMap[key]
	if !ok {
		return nil, false
	}
	if cacheEntry.expiresAt.Before(time.Now()) {
		c.deleteAndUnlink(cacheEntry)
		return nil, false
	}
	c.unlink(cacheEntry)
	c.linkFirst(cacheEntry)
	return cacheEntry, ok
}

func (c *LruCache[K, V]) TTL(key K) time.Duration {
	entry, ok := c.cacheMap[key]
	if !ok {
		return -2
	}
	if time.Now().After(entry.expiresAt) {
		return -2
	}
	return time.Until(entry.expiresAt)
}

func (c *LruCache[K, V]) deleteAndUnlink(entry *Entry[K, V]) {
	delete(c.cacheMap, entry.key)
	c.unlink(entry)
}

func (c *LruCache[K, V]) unlink(oldEntry *Entry[K, V]) {
	if oldEntry == c.head && oldEntry == c.tail {
		oldEntry.next = nil
		oldEntry.prev = nil
		c.head = nil
		c.tail = nil
		return
	} else if oldEntry == c.head {
		c.head = oldEntry.prev
		c.head.next = nil
		return
	} else if oldEntry == c.tail {
		next := oldEntry.next
		next.prev = nil
		c.tail = next
	} else {
		prev := oldEntry.prev
		next := oldEntry.next
		prev.next = next
		next.prev = prev
	}
}

func (c *LruCache[K, V]) linkFirst(entry *Entry[K, V]) {
	oldHead := c.head
	entry.prev = oldHead
	entry.next = nil
	c.head = entry
	if oldHead == nil {
		c.tail = entry
	} else {
		oldHead.next = entry
	}
}

func (c *LruCache[K, V]) setExpirationTime(entry *Entry[K, V], ttl int64) error {
	if ttl <= 0 {
		return errors.New("ttl must be greater than 0")
	}
	expiresAt := time.Now().Add(time.Duration(ttl) * time.Second)
	entry.expiresAt = expiresAt
	return nil
}

func (c *LruCache[K, V]) checkStrKey(key K) (string, K) {
	strKey, ok := any(key).(string)
	if ok {
		strKey = intern(strKey)
		return strKey, key
	}
	return "", key
}
