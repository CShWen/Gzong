package main

import (
	"net/http"
	"fmt"
	"gzong"
	"gzong/middleware"
	"encoding/base64"
)

var globalSessMgr middleware.SessionManager

func main() {
	gz := gzong.New()
	gz.GET("/testhello", helloFunc)
	gz.GET("/testjson", jsonFunc)
	gz.POST("/testpost", testPostFunc)
	//gzong.AddMiddleware(middleware.RequestDetailsLog)
	//gzong.AddMiddleware(middleware.ServiceConSumeTimeLog)
	// request headers记得添加 Authorization: [Basic c3M6cHdk]，否则请求401
	name, pwd := "ss", "pwd"
	fmt.Println("request headers记得添加 Authorization: ", "Basic "+basicAuth(name, pwd))
	u := middleware.BaseUser{Name: name, Pwd: pwd}
	gz.AddMiddleware(u.BasicAuth)

	globalSessMgr = middleware.NewSessionManager("gzCookie", 30)
	gz.GET("/login", businessAndSessionFunc)
	gz.GET("/testsm", businessAndSessionFunc)
	gz.GET("/logout", businessAndSessionFunc)

	gz.Run(":8080")
}

func helloFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello gzong.\n")
}

func jsonFunc(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"alive": true, "cc": "ss"}`))
	//io.WriteString(w, `{"alive": true, "cc": "ss"}`)
}

func testPostFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "test post.\n")
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func businessAndSessionFunc(w http.ResponseWriter, r *http.Request) {
	smResult := sessionFilter(w, r)
	if !smResult {
		fmt.Println("session验证结束，无需执行后续事务")
		return
	}
	// 后续业务代码
}

func sessionFilter(w http.ResponseWriter, r *http.Request) bool {
	sessionId, err := globalSessMgr.CheckCookieValid(w, r)
	if err != nil {
		if r.URL.Path == "/login" {
			sessionId := globalSessMgr.NewSession(w, r, make(map[interface{}]interface{}))
			fmt.Println("new SessionId:\t", sessionId)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("login success."))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
		}
		return false
	} else {
		if r.URL.Path == "/logout" {
			globalSessMgr.EndSessionById(sessionId)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("logout success."))
			return false
		} else {
			return true
		}
	}
}
