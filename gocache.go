package gocache

import (
	"fmt"
	"log"
	"sync"
)

// 定义接口 Getter
type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义函数类型 GetterFunc
type GetterFunc func(key string) ([]byte, error)

// 实现 Getter 接口的 Get 方法
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 一个 Group 是一个缓存空间，实现懂缓存的增删查方法
type Group struct {
	name      string // 缓存空间的名字
	getter    Getter // 数据源获取数据
	mainCache cache  // 主缓存，支持并发
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) // 存储所有缓存空间 Group
)

// 构造一个 Group 实例
func NewGroup(name string, capacity int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{capacity: capacity},
	}
	groups[name] = g // 把创建的 Group 实例放入全局变量 groups 中
	return g
}

// GetGroup 用来获得特定名称的 Group
func GetGroup(name string) *Group {
	mu.RLock() // 用只读锁即可，因为不涉及任何冲突变量的写操作
	g := groups[name]
	mu.RUnlock()
	return g
}

// 从缓存中根据 key 获取 value
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 从 mainCache 中查找缓存，如果存在则返回缓存值
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// 如果缓存没命中，则调用 load 方法
	return g.load(key)
}

// load 目前是直接调 getLocally
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用用户回调函数 g.getter.Get() 获取源数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	// 将源数据添加到缓存 mainCache 中
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populate 填充
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
