package main

import (
	"fmt"
	"gocache"
	"log"
	"net/http"
)

// 用map来模拟一个数据库
var db = map[string]string{
	"ShenZi":    "169.5",
	"NingGuang": "169.5",
	"KeQing":    "158.4",
}

func main() {
	groupName := "genshin_height"
	groupCapacity := int64(2 << 10)
	gocache.NewGroup(groupName, groupCapacity, gocache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}),
	)
	// 可以用 cmd 测试 curl http://localhost:9999/_gocache/genshin_height/keqing
	addr := "localhost:9999"
	peers := gocache.NewHTTPPool(addr)
	log.Println("gocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
