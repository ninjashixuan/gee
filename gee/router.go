package gee

import (
	"net/http"
)

type router struct {
	handlers map[string]HandlerHFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerHFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerHFunc) {
	key := method + "_" + pattern

	r.handlers[key] = handler

	return
}

func (r *router) handle(c *Context) {
	key := c.Method + "_" + c.Path

	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 not found: %s\n", c.Path)
	}
}