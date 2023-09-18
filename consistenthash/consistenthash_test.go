package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	// 为了测试方便，自定义一个 hash 函数，功能是把字符串转换成字符
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// 加入3个节点 被映射成虚拟节点hash 2, 4, 6, 12, 14, 16, 22, 24, 26
	// 6 -> 06 16 26
	// 4 -> 04 14 24
	// 2 -> 02 12 22
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// 再加入一个节点，映射成三个虚拟节点 8, 18, 28
	hash.Add("8")

	// 27 应该去 8号节点找，虽然节点的增删操作还没实现
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
