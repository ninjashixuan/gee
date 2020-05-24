package gee

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix string
	middleware []HandlerFunc
	parent *RouterGroup
	engine *Engine
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}

	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}

	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup{
	engine := group.engine

	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)

	return newGroup
}

func (group *RouterGroup) Use(middleware ...HandlerFunc) {
	group.middleware = append(group.middleware, middleware...)
}

func (group *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, Func HandlerFunc) {
	group.addRoute("GET", pattern, Func)
}

func (group *RouterGroup) POST(pattern string, hFunc HandlerFunc) {
	group.addRoute("POST", pattern, hFunc)
}

func (e *Engine) Run(add string) (err error) {
	return http.ListenAndServe(add, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc

	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleware...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	e.router.handle(c)
}