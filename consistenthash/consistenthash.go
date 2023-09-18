package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 函数类型，自定义实现
type Hash func(data []byte) uint32

// 一致性哈算法的主数据结构
type Map struct {
	hash     Hash           // 哈希函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环，每个虚拟节点的哈希值
	hashMap  map[int]string // 虚拟节点与真实节点的映射表，键是哈希值，值是属于哪个真实节点
}

// 构造函数，允许传入虚拟节点倍数和自定义的哈希函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	// 默认哈希函数
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 添加真实节点的函数
func (m *Map) Add(nodes ...string) {
	for _, node := range nodes {
		// 一个真实节点被映射成 replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 得到虚拟节点的哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + node)))
			// 把哈希值放进环里
			m.keys = append(m.keys, hash)
			// 记录虚拟节点到真实节点的映射
			m.hashMap[hash] = node
		}
	}
	sort.Ints(m.keys)
}

// 选择哈希节点的方法
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// 计算 key 的哈希值
	hash := int(m.hash([]byte(key)))
	// 二分查找，找到第一个使下面函数为真的节点索引
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 返回对应的真实节点
	// % 是因为上一步计算如果没找到会返回 len(m.keys)，也就是环又到第一个的情况，这时索引应该是 0
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
