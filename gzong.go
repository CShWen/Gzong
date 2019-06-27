package gzong

import (
	"context"
	"log"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

// Router web框架路由核心
type Router struct {
	srv         *http.Server
	mws         []Middleware
	handlersMap map[string]map[string]handlerFunc
}

// New 创建框架对象并将其返回
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

// Run 指定地址/端口号运行其服务
func (r *Router) Run(addr string) {
	r.srv = &http.Server{Addr: addr, Handler: r}
	err := r.srv.ListenAndServe()
	if err != nil {
		log.Fatal("error info: ", err)
	}
}

// Add 为路由添加对应指定地址、方法、处理方法
func (r *Router) Add(route string, method string, hfc handlerFunc) {
	if _, ok := r.handlersMap[route]; !ok {
		r.handlersMap[route] = make(map[string]handlerFunc)
	}
	r.handlersMap[route][method] = hfc
}

// GET 为路由添加指定地址的GET请求的处理方法
func (r *Router) GET(route string, hfc handlerFunc) {
	r.Add(route, "GET", hfc)
}

// POST 为路由添加指定地址的POST请求的处理方法
func (r *Router) POST(route string, hfc handlerFunc) {
	r.Add(route, "POST", hfc)
}

// PUT 为路由添加指定地址的PUT请求的处理方法
func (r *Router) PUT(route string, hfc handlerFunc) {
	r.Add(route, "PUT", hfc)
}

// AddMiddleware 添加中间件
func (r *Router) AddMiddleware(m Middleware) {
	r.mws = append(r.mws, m)
}

// Close 关闭服务器
func (r *Router) Close() {
	if err := r.srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Server shutdown.")
}
