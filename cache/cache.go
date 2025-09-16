package cache

import (
	"encoding/binary"
	"errors"
	"fmt"
	"ka-cache/logger"
	"strconv"
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Put(key K, value V, ttl int64) error
	Get(key K) (V, bool)
	TTL(key K) time.Duration
}

type Entry[K comparable] struct {
	key       K
	Value     []byte
	expiresAt time.Time
	next      *Entry[K]
	prev      *Entry[K]
}

type LruCache[K comparable, V any] struct {
	cacheMap    map[K]*Entry[K]
	logger      logger.Logger
	rwMutex     sync.RWMutex
	capacity    int
	cleanupStop chan bool
	head        *Entry[K]
	tail        *Entry[K]
}

func NewLruCache[K comparable, V any](cap int, logger logger.Logger) SelfClearingCache[K, V] {
	newCacheMap := make(map[K]*Entry[K], cap)
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
	strKey, key := c.conv2Str(key)
	byteVal, value := c.conv2Byte(value)
	if len(byteVal) == 0 {
		return errors.New("entry value must be string or []byte")
	}
	if len(strKey) != 0 {
		return c.putAny(any(strKey).(K), byteVal, ttl)
	}
	return c.putAny(key, byteVal, ttl)
}

func (c *LruCache[K, V]) putAny(key K, value []byte, ttl int64) error {
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
		newEntry := &Entry[K]{key: key, Value: value}
		if err := c.setExpirationTime(newEntry, ttl); err != nil {
			return err
		}
		c.cacheMap[key] = newEntry
		c.linkFirst(newEntry)
	}
	return nil
}

func (c *LruCache[K, V]) Get(key K) (V, bool) {
	var (
		entry *Entry[K]
		ok    bool
	)

	strKey, key := c.conv2Str(key)
	if len(strKey) != 0 {
		entry, ok = c.getAny(any(strKey).(K))
	} else {
		entry, ok = c.getAny(key)
	}

	var zeroV V
	if !ok {
		return zeroV, false
	}

	val, err := c.anyCast(entry.Value, key)
	if err == nil {
		return val, true
	}
	return zeroV, false
}

func (c *LruCache[K, V]) getAny(key K) (*Entry[K], bool) {
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

func (c *LruCache[K, V]) deleteAndUnlink(entry *Entry[K]) {
	delete(c.cacheMap, entry.key)
	c.unlink(entry)
}

func (c *LruCache[K, V]) unlink(oldEntry *Entry[K]) {
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

func (c *LruCache[K, V]) linkFirst(entry *Entry[K]) {
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

func (c *LruCache[K, V]) setExpirationTime(entry *Entry[K], ttl int64) error {
	if ttl <= 0 {
		return errors.New("ttl must be greater than 0")
	}
	expiresAt := time.Now().Add(time.Duration(ttl) * time.Second)
	entry.expiresAt = expiresAt
	return nil
}

func (c *LruCache[K, V]) conv2Str(key K) (string, K) {
	strKey, ok := any(key).(string)
	if ok {
		strKey = intern(strKey)
		return strKey, key
	}
	return "", key
}

func (c *LruCache[K, V]) conv2Byte(val V) ([]byte, V) {
	strVal, ok := any(val).(string)
	if ok {
		return []byte(strVal), val
	}
	return nil, val
}

func (c *LruCache[K, V]) anyCast(b []byte, key K) (out V, err error) {
	var zero V
	var anyVal any
	switch any(zero).(type) {
	case string:
		anyVal = string(b)
	case int:
		i, e := strconv.Atoi(string(b))
		if e != nil {
			return zero, e
		}
		anyVal = i
	case int64:
		i, e := strconv.ParseInt(string(b), 10, 64)
		if e != nil {
			return zero, e
		}
		anyVal = i
	case float64:
		f, e := strconv.ParseFloat(string(b), 64)
		if e != nil {
			return zero, e
		}
		anyVal = f
	case byte:
		if len(b) == 0 {
			return zero, fmt.Errorf("empty input for byte. key [%s]", key)
		}
		anyVal = b[0]
	case []byte:
		anyVal = b
	case uint32:
		if len(b) < 4 {
			return zero, fmt.Errorf("not enough bytes for uint32. key [%s]", key)
		}
		anyVal = binary.BigEndian.Uint32(b)
	default:
		return any(b).(V), nil
	}
	return anyVal.(V), nil
}
