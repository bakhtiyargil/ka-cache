package cache

import (
	"time"
)

type SelfClearingCache interface {
	StartCleanup(interval time.Duration)
	StopCleanup()
	CleanupChannel() chan bool
	Cache
}

func (c *LruCache) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.logger.Info("cache cleanup started")
			c.deleteExpiredEntries()
			c.logger.Info("cache cleanup completed")
		case <-c.CleanupChannel():
			return
		}
	}
}

func (c *LruCache) StopCleanup() {
	close(c.CleanupChannel())
}

func (c *LruCache) deleteExpiredEntries() {
	now := time.Now()
	for _, item := range c.cacheMap {
		if item.expiresAt.Before(now) {
			c.deleteAndUnlink(item)
		}
	}
}

func (c *LruCache) CleanupChannel() chan bool {
	return c.cleanupStop
}
