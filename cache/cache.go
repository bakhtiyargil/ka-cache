package cache

var SimpleCache = NewSimpleCache(128)

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

type simpleCache struct {
	cacheMap map[string]*node
	capacity int
	head     *node
	tail     *node
}

func NewSimpleCache(cap int) Cache {
	newCacheMap := make(map[string]*node, cap)
	cache := simpleCache{
		cacheMap: newCacheMap,
		capacity: cap,
		head:     nil,
		tail:     nil,
	}
	return &cache
}

func (c *simpleCache) Put(key string, value string) {
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

func (c *simpleCache) Get(key string) string {
	node, ok := c.cacheMap[key]
	if !ok {
		return ""
	}
	c.remove(node)
	c.linkFirst(node)
	return node.value
}

func (c *simpleCache) remove(node *node) {
	if node == c.head {
		return
	} else if node == c.tail {
		prev := node.prev
		prev.next = node
		node.next = nil
		c.tail = prev
	} else {
		prev := node.prev
		next := node.next
		prev.next = next
		next.prev = prev
	}
}

func (c *simpleCache) linkFirst(newNode *node) {
	oldHead := c.head
	newNode.prev = nil
	newNode.next = oldHead
	c.head = newNode
	if oldHead == nil {
		c.tail = newNode
	} else {
		oldHead.prev = newNode
	}
}
