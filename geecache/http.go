package geecache

import (
	"fmt"
	consistenthash "gee/geecache/consistenhash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self string
	basePath string

	mu sync.Mutex
	peers  *consistenthash.Map
	httpGetters  map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self,
		basePath:  defaultBasePath,
	}
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))

	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseUrl: peer + p.basePath}
	}
}

func (p *HTTPPool) PeerPicker(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("pick peer %s", peer)
		return p.httpGetters[peer], true
	}

	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %v", p.self, fmt.Sprintf(format, v))
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

type httpGetter struct {
	baseUrl string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server reeturned: %v",  res.StatusCode)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err  != nil {
		return nil, fmt.Errorf("reading responed failed: %v", err)
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)
