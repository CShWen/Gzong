package middleware

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"time"
)

// 打印请求与返回详细日志
func RequestDetailsLog(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("** request info | %s | %s | %s | %s", r.Method, r.RequestURI, r.Proto, r.RemoteAddr)
		log.Println("** request header:")
		for k, v := range r.Header {
			log.Printf("%s: %s ;", k, v)
		}
		if r.Method == "POST" || r.Method == "PUT" {
			requestBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				log.Printf("** request body: %s", requestBody)
			}
		}

		rw := httptest.NewRecorder()
		h.ServeHTTP(rw, r)

		log.Printf("** response info | %d | %s | %s ", rw.Result().StatusCode, rw.Result().Proto, rw.Result().Header)
		log.Println("** response body:", rw.Body)

		w.WriteHeader(rw.Result().StatusCode)
		for k, v := range rw.Header() {
			w.Header()[k] = v
		}
		w.Write(rw.Body.Bytes())
	}
}

// 打印服务耗时
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
