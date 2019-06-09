package Gzong

import "net/http"
// http.Handler
type Middleware func(http.HandlerFunc) http.HandlerFunc

