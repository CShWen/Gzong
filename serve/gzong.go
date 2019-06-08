package serve

import (
	"net/http"
	"log"
)

var mainRouter *Router

type handlerFunc func(http.ResponseWriter, *http.Request)

var handlersMap = make(map[string]map[string]handlerFunc)

type Router struct {
	handler http.Handler
	// todo 中间件
	mws []Middleware
}

func New() *Router {
	if mainRouter == nil {
		mainRouter = &Router{}
		mainRouter.GET("/test", testApp)
	}
	return mainRouter
}

func testApp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"test": "ok"}`))
}

func (r *Router) Run(addr string) {
	r.handler = http.HandlerFunc(routeHandler)

	err := http.ListenAndServe(addr, r.handler)
	if err != nil {
		log.Fatal("error info: ", err)
	}
}

func routeHandler(wr http.ResponseWriter, req *http.Request) {
	f, ok := handlersMap[req.URL.Path][req.Method]
	if !ok {
		http.NotFound(wr, req)
	} else {
		f(wr, req)
	}
}

func (r *Router) Add(route string, method string, hfc handlerFunc) {
	if _, ok := handlersMap[route]; !ok {
		handlersMap[route] = make(map[string]handlerFunc)
	}
	handlersMap[route][method] = hfc
}

func (r *Router) GET(route string, hfc handlerFunc) {
	r.Add(route, "GET", hfc)
}

func (r *Router) POST(route string, hfc handlerFunc) {
	r.Add(route, "POST", hfc)
}
