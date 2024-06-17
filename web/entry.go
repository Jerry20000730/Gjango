package web

import (
	"fmt"
	"github.com/Jerry20000730/Gjango/web/Constant"
	"github.com/Jerry20000730/Gjango/web/Context"
	"github.com/Jerry20000730/Gjango/web/Logic"
	"github.com/Jerry20000730/Gjango/web/Utils"
	"log"
	"net/http"
	"strconv"
)

// Handler the abstract backend logic function when router match the pattern of the URL
type Handler func(cts *context.Context)

// MiddlewareHandler the abstract backend logic function when the middleware is applied
// to the existing handle function
type MiddlewareHandler func(handler Handler) Handler

// router the struct for web router
type router struct {
	routerGroups []*routerGroup
}

// routerGroup the group of different router
// e.g., router '/user' -> handling user function
// router '/about' -> handling about page logic
type routerGroup struct {
	groupName       string
	handleFuncMap   map[string]map[string]Handler
	handleMethodMap map[string][]string
	treeNode        *Logic.TreeNode

	// for middlewares
	Middlewares   []MiddlewareHandler
	middlewareMap map[string]map[string][]MiddlewareHandler
}

// NewGroup create a new group of router
func (r *router) NewGroup(name string) *routerGroup {
	g := &routerGroup{
		groupName:       name,
		handleFuncMap:   make(map[string]map[string]Handler),
		handleMethodMap: make(map[string][]string),
		treeNode:        &Logic.TreeNode{Name: "/", Children: make([]*Logic.TreeNode, 0)},
	}
	r.routerGroups = append(r.routerGroups, g)
	return g
}

// bind function is a generic function to bind the name and the method and the handle function
func (r *routerGroup) bind(name string, method string, handler Handler) {
	if _, ok := r.handleFuncMap[name]; !ok {
		r.handleFuncMap[name] = make(map[string]Handler)
	}
	// check if the group name has already bind the designated request method
	if _, ok := r.handleFuncMap[name][method]; ok {
		panic("[ERROR] Repeated binding of request method [" + method + "] and the function")
	}
	r.handleFuncMap[name][method] = handler
	r.handleMethodMap[method] = append(r.handleMethodMap[method], name)
	r.treeNode.Put(name)
}

func (r *routerGroup) MiddlewareBind(middlewareHandler ...MiddlewareHandler) {
	r.Middlewares = append(r.Middlewares, middlewareHandler...)
}

// processHandler is a process function of how the handler actually goes, along with middlewares
func (r *routerGroup) processHandler(ctx *context.Context, handler Handler) {
	// middlewares
	if r.Middlewares != nil {
		for _, middlewareFunc := range r.Middlewares {
			handler = middlewareFunc(handler)
		}
	}
	handler(ctx)
}

// Any function allows the binding of
// 1) URL and handler
// 2) URL and request method "ANY"
func (r *routerGroup) Any(name string, handler Handler) {
	r.bind(name, Constant.ANY, handler)
}

// Get function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "GET"
func (r *routerGroup) Get(name string, handler Handler) {
	r.bind(name, Constant.GET, handler)
}

// Post function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "POST"
func (r *routerGroup) Post(name string, handler Handler) {
	r.bind(name, Constant.POST, handler)
}

// Delete function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Delete"
func (r *routerGroup) Delete(name string, handler Handler) {
	r.bind(name, Constant.DELETE, handler)
}

// Put function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Put"
func (r *routerGroup) Put(name string, handler Handler) {
	r.bind(name, Constant.PUT, handler)
}

// Patch function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Patch"
func (r *routerGroup) Patch(name string, handler Handler) {
	r.bind(name, Constant.PATCH, handler)
}

// Options function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Options"
func (r *routerGroup) Options(name string, handler Handler) {
	r.bind(name, Constant.OPTIONS, handler)
}

// Head function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Head"
func (r *routerGroup) Head(name string, handler Handler) {
	r.bind(name, Constant.HEAD, handler)
}

// Engine the main engine for web framework
type Engine struct {
	port   string
	Router router
}

// NewEngine create a new web framework engine with default port of 8321
func NewEngine() *Engine {
	return &Engine{
		port:   "8321",
		Router: router{},
	}
}

// NewEngineWithPort create a new web framework engine with user-defined port
func NewEngineWithPort(port int) *Engine {
	return &Engine{
		port:   strconv.Itoa(port),
		Router: router{},
	}
}

// GetPort get the current listening port of the web framework
func (e *Engine) GetPort() string {
	return e.port
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.httpRequestHandle(w, r)
}

func (e *Engine) httpRequestHandle(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	groups := e.Router.routerGroups
	for _, g := range groups {
		// to get the name under the group, like /hello, /get/1 so that it can be found
		// in the tree
		routerName := Utils.SubStringLast(r.RequestURI, "/"+g.groupName)
		node := g.treeNode.Get(routerName)
		if node != nil && node.IsEnd {
			// route is found
			ctx := &context.Context{
				W: w,
				R: r,
			}

			// 1. check if it is ANY method matching
			if handle, ok := g.handleFuncMap[node.Path][Constant.ANY]; ok {
				g.processHandler(ctx, handle)
				return
			}

			// 2. check if it is other method matching
			if handle, ok := g.handleFuncMap[node.Path][method]; ok {
				g.processHandler(ctx, handle)
				return
			}
			// if URL exists, but the method does not, return 405
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintf(w, "%s [%s] is not allowed\n", r.RequestURI, method)
			return
		}
	}
	// if the URL is not found, and the method, return 404
	w.WriteHeader(http.StatusNotFound)
	_, _ = fmt.Fprintf(w, "%s [%s] is not found\n", r.RequestURI, method)
}

func (e *Engine) Run() {
	//groups := e.Router.routerGroups
	//for _, g := range groups {
	//	for name, handler := range g.handleFuncMap {
	//		http.HandleFunc("/"+g.groupName+name, handler)
	//	}
	//}
	http.Handle("/", e)
	err := http.ListenAndServe(":"+e.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
