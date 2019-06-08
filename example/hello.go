package main

import (
	"net/http"
	"fmt"
	"Gzong"
)

func main() {
	gzong := Gzong.New()
	gzong.GET("/testhello", helloFunc)
	gzong.GET("/testcc", ccFunc)
	gzong.POST("/testpost", testPostFunc)
	gzong.Run(":8080")
}

func helloFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello Gzong.\n")
}

func ccFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "cc.\n")
}

func testPostFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "test post.\n")
}
