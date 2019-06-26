package gzong

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

type Router struct {
	srv         *http.Server
	mws         []Middleware
	handlersMap map[string]map[string]handlerFunc
}

func New() (r *Router) {
	r = &Router{
		handlersMap: make(map[string]map[string]handlerFunc),
	}
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f, ok := r.handlersMap[req.URL.Path][req.Method]
	if !ok {
		http.NotFound(w, req)
	} else {
		var finalHanderFunc = http.HandlerFunc(f)

		for i, l := 0, len(r.mws); i < l; i++ {
			finalHanderFunc = r.mws[i](http.HandlerFunc(finalHanderFunc))
		}
		finalHanderFunc(w, req)
	}
}

func (r *Router) Run(addr string) {
	//err := http.ListenAndServe(addr, r)
	go func() {
		r.srv = &http.Server{Addr: addr, Handler: r}
		err := r.srv.ListenAndServe()
		if err != nil {
			log.Fatal("error info: ", err)
		}
	}()
}

func (r *Router) Add(route string, method string, hfc handlerFunc) {
	if _, ok := r.handlersMap[route]; !ok {
		r.handlersMap[route] = make(map[string]handlerFunc)
	}
	r.handlersMap[route][method] = hfc
}

func (r *Router) GET(route string, hfc handlerFunc) {
	r.Add(route, "GET", hfc)
}

func (r *Router) POST(route string, hfc handlerFunc) {
	r.Add(route, "POST", hfc)
}

func (r *Router) PUT(route string, hfc handlerFunc) {
	r.Add(route, "PUT", hfc)
}

// 添加中间件
func (r *Router) AddMiddleware(m Middleware) {
	r.mws = append(r.mws, m)
}

// 关闭服务器
func (r *Router) Close() {
	if err := r.srv.Shutdown(context.Background()); err != nil {
		fmt.Println("close有错？")
		panic(err)
	}
	log.Println("Server shutdown.")
}
