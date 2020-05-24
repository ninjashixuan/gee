package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const  defaultBasePath = "/_geecache/"

type HTTPPool struct {
	self string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self,
		basePath:  defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + req.URL.Path)
	}

	p.Log("%s, %s", req.URL.Path, req.Method)

	parts := strings.SplitN(req.URL.Path[len(p.basePath):], "/", 2)

	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName :=  parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group" + groupName, http.StatusNotFound)
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

