package main

import (
	"net/http"
	"fmt"
	"Gzong"
	"Gzong/middleware"
	"encoding/base64"
)

func main() {
	gzong := Gzong.New()
	gzong.GET("/testhello", helloFunc)
	gzong.GET("/testcc", ccFunc)
	gzong.POST("/testpost", testPostFunc)
	gzong.AddMiddleware(middleware.RequestDetailsLog)
	gzong.AddMiddleware(middleware.ServiceConSumeTimeLog)
	// request headers记得添加 Authorization: [Basic c3M6cHdk]，否则请求401
	name, pwd := "ss", "pwd"
	fmt.Println("request headers记得添加 Authorization: ", "Basic "+basicAuth(name, pwd))
	u := middleware.UserInfo{name, pwd}
	gzong.AddMiddleware(u.CheckUserIdentity)
	gzong.Run(":8080")
}

func helloFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello Gzong.\n")
}

func ccFunc(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
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
