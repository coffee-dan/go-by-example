package main

import (
	"fmt"
	"net/http"
)

func main() {
	// create handler function (is this the router?)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, world")
}
