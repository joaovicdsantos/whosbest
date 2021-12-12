package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!!")
}

func main() {

	http.HandleFunc("/", helloWorld)

	log.Fatal(http.ListenAndServe(":3001", nil))
}
