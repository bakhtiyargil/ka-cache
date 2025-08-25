package cache

type Cache interface {
	Put(key string, value string)
	Get(key string) string
}

type node struct {
	key   string
	value string
	next  *node
	prev  *node
}

type LruCache struct {
	cacheMap map[string]*node
	capacity int
	head     *node
	tail     *node
}

func NewLruCache(cap int) Cache {
	newCacheMap := make(map[string]*node, cap)
	cache := LruCache{
		cacheMap: newCacheMap,
		capacity: cap,
		head:     nil,
		tail:     nil,
	}
	return &cache
}

func (c *LruCache) Put(key string, value string) {
	existingNode, ok := c.cacheMap[key]
	if ok {
		existingNode.value = value
		c.remove(existingNode)
		c.linkFirst(existingNode)
	} else {
		if len(c.cacheMap) >= c.capacity {
			delete(c.cacheMap, c.tail.key)
			c.remove(c.tail)
		}
		var newNode = node{
			key,
			value,
			nil,
			nil,
		}
		c.cacheMap[key] = &newNode
		c.linkFirst(&newNode)
	}
}

func (c *LruCache) Get(key string) string {
	node, ok := c.cacheMap[key]
	if !ok {
		return ""
	}
	c.remove(node)
	c.linkFirst(node)
	return node.value
}

func (c *LruCache) remove(oldNode *node) {
	if oldNode == c.head && oldNode == c.tail {
		oldNode.next = nil
		oldNode.prev = nil
		c.head = nil
		return
	} else if oldNode == c.head {
		c.head = oldNode.prev
		c.head.next = nil
		return
	} else if oldNode == c.tail {
		next := oldNode.next
		next.prev = nil
		c.tail = next
	} else {
		prev := oldNode.prev
		next := oldNode.next
		prev.next = next
		next.prev = prev
	}
}

func (c *LruCache) linkFirst(newNode *node) {
	oldHead := c.head
	newNode.prev = oldHead
	newNode.next = nil
	c.head = newNode
	if oldHead == nil {
		c.tail = newNode
	} else {
		oldHead.next = newNode
	}
}
