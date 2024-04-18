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

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var dbPool *pgxpool.Pool

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v\n", err)
	}

	// Set up database connection
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is not set")
	}

	dbPool, err = pgxpool.Connect(context.Background(), databaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//to close DB pool
	defer dbPool.Close()

	// Handle requests
	http.HandleFunc("/images.json", handleImages)

	// Create instance of the server
	server := &http.Server{
		Addr: ":8080",
	}

	// Start the server
	go func() {
		fmt.Fprintln(os.Stderr, "Listening on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: failed to listen and serve: %v\n", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error: failed to shutdown: %v\n", err)
	}
	log.Println("Graceful shutdown complete.")
}

// handleImages fetches images from the database and returns them as a JSON response.
func handleImages(w http.ResponseWriter, r *http.Request) {
	if dbPool == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Database connection is not set")
	}

	if r.Method == "POST" {
		postImage(w, r)
	}
	if r.Method == "GET" {
		images, err := fetchImages(dbPool, 10)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(os.Stderr, "Error: failed to fetch images: %v\n", err)
		}

		indent, err := getIndentParam(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(os.Stderr, "Error: failed to parse indent: %v\n", err)
		}

		b, err := json.MarshalIndent(images, "", indent)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(os.Stderr, "Error: failed to marshal images: %v\n", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// POST /images.json
func postImage(w http.ResponseWriter, r *http.Request) {

	// Parse the request body
	var image types.Image
	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Error: failed to decode image: %v\n", err)
		return
	}

	// Check if the image url exists in the database
	query := "SELECT url FROM public.images WHERE url = $1"
	row := dbPool.QueryRow(context.Background(), query, image.URL)
	var url string
	err = row.Scan(&url)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Image with url already exists\n"))
		fmt.Fprintf(os.Stderr, "Error: image with url %s already exists\n", image.URL)
		return
	}

	// Insert the image into the database
	query = "INSERT INTO public.images (title, url, alt_text) VALUES ($1, $2, $3)"
	_, err = dbPool.Exec(context.Background(), query, image.Title, image.URL, image.AltText)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Error: failed to insert image: %v\n", err)
	}
	w.WriteHeader(http.StatusCreated)
}

// fetchImages fetches images from the database.
func fetchImages(pool *pgxpool.Pool, limit int) ([]types.Image, error) {
	var images []types.Image
	query := fmt.Sprintf("SELECT title, url, alt_text FROM public.images LIMIT %d", limit)
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to fetch images: %v\n", err)
	}
	for rows.Next() {
		var title, url, altText string
		err = rows.Scan(&title, &url, &altText)
		images = append(images, types.Image{Title: title, URL: url, AltText: altText})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to scan image: %v\n", err)
		}
	}
	return images, err
}

// getIndentParam returns the indent parameter from the request query string.
func getIndentParam(r *http.Request) (string, error) {
	params := r.URL.Query()
	indent := params.Get("indent")
	// case when indent is not provided
	if indent == "" {
		return "", nil
	}
	// case when indent is not parsable to int
	if condition, err := strconv.ParseInt(indent, 10, 8); condition == 0 {
		return "", err
	}

	indentSize, err := strconv.Atoi(indent)
	if err != nil {
		return "", err
	}
	indent = ""
	for i := 0; i < indentSize; i++ {
		indent += " "
	}
	return indent, nil
}
