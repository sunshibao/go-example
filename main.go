package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "哎哟不错哦", "3333333")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
