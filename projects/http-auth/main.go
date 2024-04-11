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

// http routes
var routes = map[string]http.HandlerFunc{
	"/":              rootHandler,
	"/200":           successHandler,
	"/404":           http.NotFound,
	"/500":           serverErrorHandler,
	"/authenticated": authenticatedHandler,
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file : %v", err)
	}

	for route, handler := range routes {
		if route == "/" || route == "/authenticated" {
			http.HandleFunc(route, rateLimiterMiddleware(handler, 100, 30))
		} else {
			http.HandleFunc(route, handler)
		}
	}

	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}
}

func rateLimiterMiddleware(next http.HandlerFunc, rl rate.Limit, b int) http.HandlerFunc {
	limiter := rate.NewLimiter(rl, b)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}
	})
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("200\n"))
}

func serverErrorHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain")
	w.WriteHeader(500)
	w.Write([]byte("500\n"))
}

func authenticatedHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	USERNAME := os.Getenv("AUTH_USERNAME")
	PASSWORD := os.Getenv("AUTH_PASSWORD")

	if !ok || username != USERNAME || password != PASSWORD {
		w.Header().Set("WWW-Authenticate", `Basic realm="localhost"`)
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
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
}
