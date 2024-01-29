package main

type node struct {
	key   string
	value string
	next  *node
	prev  *node
}

type Cache struct {
	cacheMap map[string]*node
	capacity int
	head     *node
	tail     *node
}

func (c *Cache) Set(key string, value string) {
	existingNode, ok := c.cacheMap[key]
	if ok {
		existingNode.value = value
		c.remove(existingNode)
		c.setHead(existingNode)
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
		c.setHead(&newNode)
	}

}

func (c *Cache) Get(key string) string {
	node, ok := c.cacheMap[key]
	if !ok {
		return ""
	}
	c.remove(node)
	c.setHead(node)
	return node.value
}

func (c *Cache) remove(node *node) *node {
	predecessor := node.prev
	successor := node.next
	if c.head != node {
		predecessor.next = successor
		successor.prev = predecessor
	}
	return node
}

func (c *Cache) setHead(node *node) {
	node.next = c.head.next
	node.prev = c.head
	c.head.next = node
	c.head = node
}

func NewCache(cap int) *Cache {
	newHead := node{
		key:   "",
		value: "",
		next:  nil,
		prev:  nil,
	}
	newTail := node{
		key:   "",
		value: "",
		next:  &newHead,
		prev:  nil,
	}
	newCacheMap := make(map[string]*node, cap)

	cache := Cache{
		cacheMap: newCacheMap,
		capacity: cap,
		head:     &newHead,
		tail:     &newTail,
	}
	return &cache
}
