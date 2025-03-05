package main

import (
	"fmt"
	"net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "Hello, this is an echo server!")
	if err != nil {
		panic("bad things happened")
	}
}

func main() {
	http.HandleFunc("/", echoHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
