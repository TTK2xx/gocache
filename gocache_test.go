package gocache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

// 用map来模拟一个数据库
var db = map[string]string{
	"ShenZi":    "169.5",
	"NingGuang": "169.5",
	"KeQing":    "158.4",
}

// 借助 GetterFunc 的类型转换，将一个匿名回调函数转换成了接口 f Getter
func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed")
	}
}

func TestGet(t *testing.T) {
	// 统计数据从数据库中被读了几次
	loadCounts := make(map[string]int, len(db))
	groupName := "genshin_height"
	groupCapacity := int64(2 << 10)
	// 初始化一个group (名称， 容量， 从数据源获取数据的回调函数)
	group := NewGroup(groupName, groupCapacity, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		// 第一次会从 DB 中加载，测试是否加载正确了
		if view, err := group.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		// 从 DB 中加载一次就会在缓存里取了，所以 Count 应该是1
		if _, err := group.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	// 测试取不存在的值
	if view, err := group.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

func TestGetGroup(t *testing.T) {
	groupName := "genshin_height"
	groupCapacity := int64(2 << 10)
	NewGroup(groupName, groupCapacity, GetterFunc(
		func(key string) (bytes []byte, err error) { return }))
	// 测试获取存在的 group
	if group := GetGroup(groupName); group == nil || group.name != groupName {
		t.Fatalf("group %s not exist", groupName)
	}
	// 测试获取不存在的 group
	if group := GetGroup(groupName + "111"); group != nil {
		t.Fatalf("expect nil, but %s got", group.name)
	}
}
