package lru // 最近最少使用缓存淘汰策略，缓存满时删除最久没有使用的数据。

import "container/list"

// Cache 是lru缓存的主要结构体
type Cache struct {
	maxBytes  int64                         // 最大允许使用内存
	nbytes    int64                         // 当前已经使用内存
	ll        *list.List                    // 双向链表，维护访问顺序
	cache     map[string]*list.Element      // 哈希表，快速查找
	onEvicted func(key string, value Value) // 删除回调函数
}

// entry 是链表中存储的数据结构
type entry struct {
	key   string
	value Value
}

// Value 接口，所有存储的值需要实现这个接口
type Value interface {
	Len() int // 返回值占用的内存大小
}

// New 是一个构造函数，用来创建一个新的 Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Add 方法用来向缓存中增加一个元素
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get 方法用来查找键对应的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 缓存淘汰，移除最近最少访问的节点，队首
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

// Len 返回缓存元素个数
func (c *Cache) Len() int {
	return c.ll.Len()
}
