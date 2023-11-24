package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// create handler function (is this the router?)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if content, error := os.ReadFile("views/index.html"); error != nil {
		log.Fatal(error)
	} else {
		// fmt.Fprint(w, string(content))
		w.Write(content)
	}
}
