package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", hello)
	fmt.Printf("listening on localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("/hello received request\n")
	fmt.Fprintf(w, "Hello World!")
}