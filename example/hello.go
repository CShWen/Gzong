package main

import (
	"net/http"
	"fmt"
	"Gzong"
	"Gzong/middleware"
)

func main() {
	gzong := Gzong.New()
	gzong.GET("/testhello", helloFunc)
	gzong.GET("/testcc", ccFunc)
	gzong.POST("/testpost", testPostFunc)
	//gzong.AddMiddleware(middleware.Logtest)
	//gzong.AddMiddleware(middleware.Logtest2)
	gzong.AddMiddleware(middleware.RequestDetailsLog)
	gzong.AddMiddleware(middleware.ServiceConSumeTimeLog)
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
