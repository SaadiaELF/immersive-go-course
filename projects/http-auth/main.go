package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file : %v", err)
	}
	limiter := rate.NewLimiter(100, 30)
	if limiter.Allow() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			params := r.URL.Query()
			w.Header().Add("Content-Type", "text/html")

			if r.Method == "POST" {
				body, err := io.ReadAll(r.Body)

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
				w.Write(body)
			}

			if r.Method == "GET" {
				if len(params) > 0 {
					for key, values := range params {
						values := html.EscapeString((values[0]))
						strings := fmt.Sprintf("<!DOCTYPE html><html><em>Hello, world</em><p>Query parameters:<ul><li>%v:%v</li></ul>\n", key, values)
						w.Write([]byte(strings))
					}
				} else {
					w.Write([]byte("<!DOCTYPE html><html><em>Hello, world</em>\n"))
				}
			}

			fmt.Fprintln(os.Stderr, "Error: invalid request")

		})

		// authentication
		http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			USERNAME := os.Getenv("AUTH_USERNAME")
			PASSWORD := os.Getenv("AUTH_PASSWORD")

			if !ok || username != USERNAME || password != PASSWORD {
				w.Header().Set("WWW-Authenticate", `Basic realm="localhost"`)
				w.WriteHeader(http.StatusUnauthorized)
			}
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

	} else {
		fmt.Fprintln(os.Stderr, "Error: rate limit exceeded")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}
}
