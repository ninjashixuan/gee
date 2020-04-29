package gee

import "net/http"

type HandlerHFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerHFunc) {
	e.router.addRoute(method, pattern, handler)
}

func (e *Engine) Get(pattern string, Func HandlerHFunc) {
	e.router.addRoute("GET", pattern, Func)
}

func (e *Engine) POST(pattern string, hFunc HandlerHFunc) {
	e.router.addRoute("POST", pattern, hFunc)
}

func (e *Engine) Run(add string) (err error) {
	return http.ListenAndServe(add, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}