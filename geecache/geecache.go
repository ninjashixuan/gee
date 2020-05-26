package geecache

import (
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache cache

	peers PeerPicker
}

var (
	m      sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter nil")
	}

	m.Lock()
	defer m.Unlock()

	group := &Group{
		name:      name,
		mainCache: cache{cacheBytes: cacheBytes},
		getter:    getter,
	}

	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	m.RLock()
	g := groups[name]
	m.RUnlock()

	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than one")
	}

	g.peers = peers
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[geecache hit]")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peers, ok := g.peers.PeerPicker(key); ok {
			if value, err = g.getFromPeer(peers, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer", err)

		}
	}
	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: bytes}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Add(key, value)
}

type Getter interface {
	Get(string) ([]byte, error)
}

type GetterFunc func(string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
