package gzong

import "net/http"

// Middleware 自定义的中间件方法，入参出参均为http.HandlerFunc
type Middleware func(http.HandlerFunc) http.HandlerFunc
