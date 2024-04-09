package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world\n"))
	})
	// 200 status code
	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("200\n"))
	})

	// 500 status code
	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("Internal server error\n"))
	})

	// 404 status code
	http.Handle("/404", http.NotFoundHandler())

	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}

}
