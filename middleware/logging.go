package middleware

import (
	"log"
	"net/http"
	"fmt"
)

//func Logtest(h http.HandlerFunc) http.HandlerFunc {
func Logtest(h http.HandlerFunc) http.HandlerFunc {
	log.Println("testFuncLog.")
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("testFuncLog_1")
		h.ServeHTTP(w, r)
		fmt.Println("testFuncLog_2")
	}
}

func Logtest2(h http.HandlerFunc) http.HandlerFunc {
	log.Println("测试日志.")
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("测试日志Func_1")
		h.ServeHTTP(w, r)
		fmt.Println("测试日志Func_2")
	}
}
