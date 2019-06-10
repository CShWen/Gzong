package middleware

import (
	"log"
	"net/http"
	"fmt"
	"time"
	"net/http/httptest"
	"runtime"
	"io/ioutil"
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

func RequestDetailsLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("** request info | %s | %s | %s | %s", r.Method, r.RequestURI, r.Proto, r.RemoteAddr)
		log.Println("** request header:")
		for k, v := range r.Header {
			log.Printf("%s: %s ;", k, v)
		}
		if r.Method == "POST" {
			requestBody, ok := ioutil.ReadAll(r.Body)
			if ok == nil {
				log.Printf("** request body: %s", requestBody)
			}
		}

		newRw := httptest.NewRecorder()
		h.ServeHTTP(newRw, r)
		for k, v := range newRw.Header() {
			w.Header()[k] = v
		}
		//w.Header()["testHeader"] = []string{"t1", "t2"}
		//w.WriteHeader(newRw.Code)
		//w.Write(newRw.Body.Bytes())

		log.Printf("** response info | %d | %s | %s ", newRw.Result().StatusCode, newRw.Result().Proto, newRw.Result().Header)
		log.Println("** response body:", newRw.Body)
	}
}

func ServiceConSumeTimeLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now().UnixNano()
		h.ServeHTTP(w, r)
		t2 := time.Now().UnixNano()
		ct := float64(t2-t1) / 1e6

		pc, _, _, _ := runtime.Caller(0)
		log.Printf("%s | %s | %s | %.3fms \n", runtime.FuncForPC(pc).Name(), r.Method, r.URL.Path, ct)
	}
}
