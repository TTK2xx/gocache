package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

// 测试Get方法
func TestGet(t *testing.T) {
	lru := NewCache(int64(128), nil)
	lru.Add("k1", String("v1"))
	// 测试读存在的值
	if v, ok := lru.Get("k1"); !ok || string(v.(String)) != "v1" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	// 测试读不存在的值
	if _, ok := lru.Get("k2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

// 测试 当使用内存超过了设定值时，是否会触发正确的节点的移除
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("k1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest k1 failed")
	}
}

// 测试回调函数和Get()
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	// 假设回调函数会把淘汰的 key 加到 keys 中
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := NewCache(int64(10), callback)
	lru.Add("k1", String("v1"))
	lru.Add("k2", String("v2"))
	lru.Get("k1") // 通过获k1来刷新他的最近使用，防止被淘汰
	lru.Add("k3", String("v3"))
	lru.Get("k1")
	lru.Add("k4", String("v4"))

	// 触发淘汰策略，会把前两个删除
	// expect := []string{"k1", "k2"}
	expect := []string{"k2", "k3"}
	// expect := []string{"k1", "k2", "k3"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
