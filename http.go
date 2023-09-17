package gocache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_gocache/"

type HTTPPool struct {
	// // 记录一个节点自己的地址, e.g. "https://example.net/_gocache/"
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 打印日志的方法
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 处理 HTTP 请求的方法
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 首先判断访问路径的前缀是否是 basePath，不是返回错误。
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// 约定格式： /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将缓存值作为 httpResponse 的 body 返回。
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
