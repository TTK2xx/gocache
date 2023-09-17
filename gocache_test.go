package gocache

import (
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
