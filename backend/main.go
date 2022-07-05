package main

import (
	"fmt"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Web!\n"))
}

func main() {
	fmt.Println("Hello, world!")
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	addr := ":8084"
	http.ListenAndServe(addr, mux)

}
