package web

import (
	"fmt"
	"github.com/Jerry20000730/Gjango/web/Constant"
	"github.com/Jerry20000730/Gjango/web/Context"
	"github.com/Jerry20000730/Gjango/web/File"
	"github.com/Jerry20000730/Gjango/web/Logic"
	"github.com/Jerry20000730/Gjango/web/Render"
	"github.com/Jerry20000730/Gjango/web/Utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
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
		Middlewares:     make([]MiddlewareHandler, 0),
		middlewareMap:   make(map[string]map[string][]MiddlewareHandler),
	}
	r.routerGroups = append(r.routerGroups, g)
	return g
}

// bind function is a generic function to bind the name and the method and the handle function
// the bind function support bind function to specific router group and request method
// it also bind middleware functions to specific router group and specific request method
func (r *routerGroup) bind(name string, method string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	if _, ok := r.handleFuncMap[name]; !ok {
		r.handleFuncMap[name] = make(map[string]Handler)
		r.middlewareMap[name] = make(map[string][]MiddlewareHandler)
	}
	// check if the group name has already bind the designated request method
	if _, ok := r.handleFuncMap[name][method]; ok {
		panic("[ERROR] Repeated binding of request method [" + method + "] and the function")
	}
	r.handleFuncMap[name][method] = handler
	r.middlewareMap[name][method] = append(r.middlewareMap[name][method], middlewareHandler...)
	r.handleMethodMap[method] = append(r.handleMethodMap[method], name)
	r.treeNode.Put(name)
}

func (r *routerGroup) MiddlewareRegister(middlewareHandler ...MiddlewareHandler) {
	r.Middlewares = append(r.Middlewares, middlewareHandler...)
}

// processHandler is a process function of how the handler actually goes, along with middlewares
func (r *routerGroup) processHandler(name string, method string, ctx *context.Context, handler Handler) {
	// router-group general middlewares
	if r.Middlewares != nil {
		for _, middlewareFunc := range r.Middlewares {
			handler = middlewareFunc(handler)
		}
	}

	// router-group request method specific middlewares
	routerMiddlewares := r.middlewareMap[name][method]
	if routerMiddlewares != nil {
		for _, middlewareFunc := range routerMiddlewares {
			handler = middlewareFunc(handler)
		}
	}

	handler(ctx)
}

// Any function allows the binding of
// 1) URL and handler
// 2) URL and request method "ANY"
func (r *routerGroup) Any(name string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	r.bind(name, Constant.ANY, handler, middlewareHandler...)
}

// Get function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "GET"
func (r *routerGroup) Get(name string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	r.bind(name, Constant.GET, handler, middlewareHandler...)
}

// Post function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "POST"
func (r *routerGroup) Post(name string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	r.bind(name, Constant.POST, handler, middlewareHandler...)
}

// Delete function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Delete"
func (r *routerGroup) Delete(name string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	r.bind(name, Constant.DELETE, handler, middlewareHandler...)
}

// Put function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Put"
func (r *routerGroup) Put(name string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	r.bind(name, Constant.PUT, handler, middlewareHandler...)
}

// Patch function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Patch"
func (r *routerGroup) Patch(name string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	r.bind(name, Constant.PATCH, handler, middlewareHandler...)
}

// Options function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Options"
func (r *routerGroup) Options(name string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	r.bind(name, Constant.OPTIONS, handler, middlewareHandler...)
}

// Head function allows the binding of
// 1) URL and handler
// 2) URL and HTTP request method "Head"
func (r *routerGroup) Head(name string, handler Handler, middlewareHandler ...MiddlewareHandler) {
	r.bind(name, Constant.HEAD, handler, middlewareHandler...)
}

// Engine the main engine for web framework
type Engine struct {
	port          string
	pool          sync.Pool
	Router        router
	HTMLPreloader Render.HTMLPreloader
	FileManager   File.FileManager
}

// NewEngine create a new web framework engine with default port of 8321
func NewEngine() *Engine {
	engine := &Engine{
		port:   "8321",
		Router: router{},
	}
	engine.pool.New = func() any {
		return &context.Context{}
	}
	return engine
}

// NewEngineWithPort create a new web framework engine with user-defined port
func NewEngineWithPort(port int) *Engine {
	engine := &Engine{
		port:   strconv.Itoa(port),
		Router: router{},
	}
	engine.pool.New = func() any {
		return &context.Context{}
	}
	return engine
}

// GetPort get the current listening port of the web framework
func (e *Engine) GetPort() string {
	return e.port
}

// ServeHTTP the function that process http request
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*context.Context)
	ctx.W = w
	ctx.R = r
	e.httpRequestHandle(ctx, w, r)
	e.pool.Put(ctx)
}

// PreLoadFuncMap the function that pre-read the template.funcmap into the memory
func (e *Engine) PreLoadFuncMap(funcMap template.FuncMap) {
	e.HTMLPreloader.FuncMap = funcMap
}

// PreLoadTemplate the function that pre-read the template (html files) into the memory
func (e *Engine) PreLoadTemplate(pattern string) {
	t := template.Must(template.New("").Funcs(e.HTMLPreloader.FuncMap).ParseGlob(pattern))
	e.HTMLPreloader.Template = t
}

func (e *Engine) httpRequestHandle(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	method := r.Method
	groups := e.Router.routerGroups
	for _, g := range groups {
		// to get the name under the group, like /hello, /get/1 so that it can be found
		// in the tree
		routerName := Utils.SubStringLast(r.RequestURI, "/"+g.groupName)
		node := g.treeNode.Get(routerName)
		if node != nil && node.IsEnd {

			// 1. check if it is ANY method matching
			if handle, ok := g.handleFuncMap[node.Path][Constant.ANY]; ok {
				g.processHandler(routerName, Constant.ANY, ctx, handle)
				return
			}

			// 2. check if it is other method matching
			if handle, ok := g.handleFuncMap[node.Path][method]; ok {
				g.processHandler(routerName, method, ctx, handle)
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
