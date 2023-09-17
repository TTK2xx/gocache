package lru

import "container/list"

// LRU缓存
type Cache struct {
	capacity  int64                         // 允许使用的最大缓存容量，0代表无限制
	length    int64                         // 当前已使用的缓存容量
	ll        *list.List                    // 一个双向链表，用于存实际值和维护最近最少使用的顺序
	cache     map[string]*list.Element      // 键是字符串，值是双向链表中对应节点的指针
	OnEvicted func(key string, value Value) // 某条记录被移除的时候的回调函数
}

// 键值对 entry 是双向链表节点的数据类型，在链表中冗余每个值对应的 key 的目的在于，
// 淘汰队首节点时，需要用 key 从 map 中删除对应的映射。
type entry struct {
	key   string
	value Value
}

// 为了通用性，允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法
// Len() int， 用于返回值所占用的内存大小。
type Value interface {
	Len() int
}

// Cache的构造函数
func NewCache(capacity int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		capacity:  capacity,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查找 key 对应的 value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 将链表中的节点 ele 移动到队尾
		c.ll.MoveToBack(ele)
		// .(*entry) 是一个类型断言，它的目的是将 ele.Value 的值转换为类型为 entry 的指针。
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// 缓存淘汰，即移除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Front() // 取到队首节点
	if ele != nil {
		c.ll.Remove(ele) // 从链表中删除
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // 从map中删除
		c.length -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新当前所用内存的大小
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 新增/修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 如果存在，即是新增
		c.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		c.length += int64(value.Len()) - int64(kv.value.Len()) // 更新当前所用内存的大小
		kv.value = value
	} else {
		// 如果不存在，即是修改
		ele := c.ll.PushBack(&entry{key, value})
		c.cache[key] = ele
		c.length += int64(len(key)) + int64(value.Len())
	}
	// 如果超过最大容量，循环移除最近最少使用的缓存项
	for c.capacity != 0 && c.capacity < c.length {
		c.RemoveOldest()
	}
}

// 缓存的数量
func (c *Cache) Len() int {
	return c.ll.Len()
}
