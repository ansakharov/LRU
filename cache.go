package lru

type Cache struct {
	cache    map[string]*dLinkedNode
	size     int
	capacity int
	headNode *dLinkedNode
	tailNode *dLinkedNode
}

type dLinkedNode struct {
	key   string
	value interface{}

	prev *dLinkedNode
	next *dLinkedNode
}

func New(capacity uint) *Cache {
	head := &dLinkedNode{}
	tail := &dLinkedNode{}

	head.next = tail
	tail.prev = head

	return &Cache{
		cache:    make(map[string]*dLinkedNode, capacity+1),
		size:     0,
		capacity: int(capacity),
		headNode: head,
		tailNode: tail,
	}
}

func (c *Cache) Get(key string) interface{} {
	node, ok := c.cache[key]
	if !ok {
		return nil
	}
	c.moveToHead(node)

	return node.value
}

func (c *Cache) Set(key string, value interface{}) {
	node, ok := c.cache[key]
	if ok {
		node.value = value
		c.moveToHead(node)
	} else {
		newNode := &dLinkedNode{
			value: value,
			key:   key,
		}

		c.cache[key] = newNode
		c.addNode(newNode)

		c.size++

		if c.size > c.capacity {
			obsoleteCache := c.popTail()
			delete(c.cache, obsoleteCache.key)
			c.size--
		}
	}
}

func (c *Cache) addNode(node *dLinkedNode) {
	node.prev = c.headNode
	node.next = c.headNode.next

	c.headNode.next.prev = node
	c.headNode.next = node
}

func (c *Cache) removeNode(node *dLinkedNode) {
	node.next.prev = node.prev
	node.prev.next = node.next
}

func (c *Cache) moveToHead(node *dLinkedNode) {
	c.removeNode(node)
	c.addNode(node)
}

func (c *Cache) popTail() *dLinkedNode {
	node := c.tailNode.prev
	c.removeNode(node)

	return node
}
