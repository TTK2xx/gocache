package gocache

// 存储真实缓存值的数据结构，byte数组更通用
type ByteView struct {
	b []byte
}

// 实现 Len() int 方法，定义在 lru 中，Value 必须实现这个方法
func (v ByteView) Len() int {
	return len(v.b)
}

// 使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改。
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v ByteView) String() string {
	return string(v.b)
}
