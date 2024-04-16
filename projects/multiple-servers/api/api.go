package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"servers/images"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbPool *pgxpool.Pool
var err error

func Run(dbURL *string, port *string) {
	dbPool, err = pgxpool.Connect(context.Background(), *dbURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//closing database connections
	defer dbPool.Close()

	// Handle requests
	http.HandleFunc("/api/images.json", handleImages)

	// Create instance of http.Server
	server := &http.Server{
		Addr: ":" + *port,
	}

	// Start the server in a goroutine
	go func() {
		fmt.Fprintln(os.Stderr, "Listening on port", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintln(os.Stderr, "Error starting server:", err)
		}
		fmt.Println("Server stopped serving new requests.")
	}()

	// Wait for a signal to shut down the server
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Shut down the server
	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "Error shutting down server:", err)
	}
	fmt.Println("Server shut down.")
}

// handleImages fetches images from the database and returns them as a JSON response.
func handleImages(w http.ResponseWriter, r *http.Request) {
	if dbPool == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Database connection is not set")
	}

	if r.Method == "POST" {
		images.PostImage(dbPool, w, r)
	}
	if r.Method == "GET" {
		img, err := images.FetchImages(dbPool, 10)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(os.Stderr, "Error: failed to fetch images: %v\n", err)
		}

		indent, err := images.GetIndentParam(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(os.Stderr, "Error: failed to parse indent: %v\n", err)
		}

		b, err := json.MarshalIndent(img, "", indent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(os.Stderr, "Error: failed to marshal images: %v\n", err)
		}

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082")
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
