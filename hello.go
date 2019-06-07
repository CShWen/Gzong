package main

import (
	"net/http"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("/", helloFunc)
	err := http.ListenAndServe(":8080", nil)
	if (err != nil) {
		log.Fatal("error info: ", err)
	}
}
func helloFunc(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "hello Gzong.\n")
}
