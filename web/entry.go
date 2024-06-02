package web

import (
	"fmt"
	"github.com/Jerry20000730/Gjango/web/Constant"
	"github.com/Jerry20000730/Gjango/web/Context"
	"log"
	"net/http"
	"strconv"
)

// Handler the abstract backend logic function when router match the pattern of the URL
type Handler func(cts *context.Context)

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
}

// NewGroup create a new group of router
func (r *router) NewGroup(name string) *routerGroup {
	g := &routerGroup{
		groupName:       name,
		handleFuncMap:   make(map[string]map[string]Handler),
		handleMethodMap: make(map[string][]string),
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
	method := r.Method
	groups := e.Router.routerGroups
	for _, g := range groups {
		for name, methodHandlerIndex := range g.handleFuncMap {
			url := "/" + g.groupName + name
			// if the router is matched
			if r.RequestURI == url {
				ctx := &context.Context{
					W: w,
					R: r,
				}

				// 1. check if it is ANY method matching
				if handle, ok := methodHandlerIndex[Constant.ANY]; ok {
					handle(ctx)
					return
				}

				// 2. check if it is other method matching
				if handle, ok := methodHandlerIndex[method]; ok {
					handle(ctx)
					return
				}
				// if URL exists, but the method does not, return 405
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, "%s [%s] is not allowed\n", r.RequestURI, method)
				return
			}
		}
	}
	// if the URL is not found, and the method, return 404
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s [%s] is not found\n", r.RequestURI, method)
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