package cache

import (
	"errors"
	"ka-cache/logger"
	"sync"
	"time"
)

type Cache interface {
	Put(key string, value string, ttl int64) error
	Get(key string) (*Entry, bool)
	TTL(key string) time.Duration
}

type Entry struct {
	key       string
	Value     string
	expiresAt time.Time
	next      *Entry
	prev      *Entry
}

type LruCache struct {
	cacheMap    map[string]*Entry
	logger      logger.Logger
	rwMutex     sync.RWMutex
	capacity    int
	cleanupStop chan bool
	head        *Entry
	tail        *Entry
}

func NewLruCache(cap int, logger logger.Logger) SelfClearingCache {
	newCacheMap := make(map[string]*Entry, cap)
	cache := LruCache{
		cacheMap: newCacheMap,
		capacity: cap,
		logger:   logger,
		head:     nil,
		tail:     nil,
	}
	return &cache
}

func (c *LruCache) Put(key string, value string, ttl int64) error {
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
		var newEntry Entry
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

func (c *LruCache) Get(key string) (*Entry, bool) {
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

func (c *LruCache) TTL(key string) time.Duration {
	entry, ok := c.cacheMap[key]
	if !ok {
		return -2
	}
	if time.Now().After(entry.expiresAt) {
		return -2
	}
	return time.Until(entry.expiresAt)
}

func (c *LruCache) deleteAndUnlink(entry *Entry) {
	delete(c.cacheMap, entry.key)
	c.unlink(entry)
}

func (c *LruCache) unlink(oldEntry *Entry) {
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

func (c *LruCache) linkFirst(entry *Entry) {
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

func (c *LruCache) setExpirationTime(entry *Entry, ttl int64) error {
	if ttl <= 0 {
		return errors.New("ttl must be greater than 0")
	}
	expiresAt := time.Now().Add(time.Duration(ttl) * time.Second)
	entry.expiresAt = expiresAt
	return nil
}
