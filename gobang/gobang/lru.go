package gobang

type linkNode struct {
	key  MenualKey
	prev *linkNode
	next *linkNode
}
type cacheNode struct {
	val   *MenualItem
	lnode *linkNode
}

type LRUCache struct {
	cache    map[MenualKey]cacheNode
	head     *linkNode
	tail     *linkNode
	Capacity int
	Used     int
}

func NewLRUCache(capacity int) *LRUCache {
	var lru LRUCache
	lru.Capacity = capacity
	lru.cache = make(map[MenualKey]cacheNode)
	return &lru
}

func (this *LRUCache) promote(lnode *linkNode) {
	if lnode == this.head {
		return
	}
	if lnode.prev != nil {
		if this.tail == lnode {
			this.tail = lnode.prev
		}
		lnode.prev.next = lnode.next
	}
	if lnode.next != nil {
		lnode.next.prev = lnode.prev
	}
	lnode.next = this.head
	if this.head != nil {
		this.head.prev = lnode
	}
	lnode.prev = nil
	this.head = lnode
}

func (this *LRUCache) Get(key MenualKey) *MenualItem {
	cnode, ok := this.cache[key]
	if !ok || cnode.val == nil {
		return nil
	}
	//promote the linknode to head
	this.promote(cnode.lnode)
	return cnode.val
}

func (this *LRUCache) Put(key MenualKey, value *MenualItem) {
	cnode, ok := this.cache[key]
	if !ok || cnode.val == nil {
		// if exceed the capacity delete from tail
		if this.Used == this.Capacity {
			key1 := this.tail.key
			this.tail = this.tail.prev
			if this.tail != nil {
				this.tail.next = nil
			}
			delete(this.cache, key1)
			//this.cache[key1] = cacheNode{nil, nil}
			//fmt.Println("delete", key1)
		} else {
			this.Used++
		}
		// create new link node
		lnode := new(linkNode)
		lnode.key = key
		// append the new node to head
		lnode.next = this.head
		if this.Used != 1 {
			this.head.prev = lnode
		} else {
			// for the first node, point tail to it
			this.tail = lnode
		}
		this.head = lnode

		// append {value, *node} to cache
		this.cache[key] = cacheNode{value, lnode}
	} else {
		// promote the linknode to head
		this.promote(cnode.lnode)
		// set value to cache
		this.cache[key] = cacheNode{value, cnode.lnode}
	}
}
