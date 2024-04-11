package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server-database/types"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var dbPool *pgxpool.Pool

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	// Set up database connection
	dbPool, err := setDatabaseConnection()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	defer dbPool.Close()

	// Create instance of the server
	server := &http.Server{
		Addr: ":8080",
	}

	// Handle requests
	http.HandleFunc("/images.json", handleImages)

	// Start the server
	go func() {
		fmt.Fprintln(os.Stderr, "Listening on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: failed to listen and serve: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error: failed to shutdown: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}

func setDatabaseConnection() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is not set")
	}

	pool, err := pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to connect to database: %v", err)
	}
	return pool, err
}

func handleImages(w http.ResponseWriter, r *http.Request) {
	images, err := fetchImages(dbPool)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	indent := getIndentParam(r)
	b, err := json.MarshalIndent(images, "", indent)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func fetchImages(pool *pgxpool.Pool) ([]types.Image, error) {
	var images []types.Image
	rows, err := pool.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to fetch images: %v", err)
	}
	for rows.Next() {
		var title, url, altText string
		err = rows.Scan(&title, &url, &altText)
		images = append(images, types.Image{Title: title, URL: url, AltText: altText})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to scan image: %v", err)
		}
	}
	return images, err
}

func getIndentParam(r *http.Request) string {
	params := r.URL.Query()

	indentSize, err := strconv.Atoi(params.Get("indent"))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	indent := ""
	for i := 0; i < indentSize; i++ {
		indent += " "
	}
	return indent
}
