package main

import (
	"fmt"
	"github.com/cshwen/gzong"
	"net/http"
)

func main() {
	gz := gzong.New()
	gz.GET("/test", testFunc)
	gz.Run(":8080")
}

func testFunc(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "hello gzong.")
}
