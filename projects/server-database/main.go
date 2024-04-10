package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server-database/types"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func handleImages(w http.ResponseWriter, r *http.Request) {
	images, _ := fetchImages()
	indent := getIndentParam(r)
	b, err := json.MarshalIndent(images, "", indent)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
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

func fetchImages() ([]types.Image, error) {
	images := []types.Image{
		{
			Title:   "Sunset",
			AltText: "Clouds at sunset",
			URL:     "https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
		},
		{
			Title:   "Mountain",
			AltText: "A mountain at sunset",
			URL:     "https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80",
		},
	}

	return images, nil
}
func setDatabaseConnection() (*pgx.Conn, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is not set")
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to connect to database: %v", err)
	}
	return conn, err
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	setDatabaseConnection()
	http.HandleFunc("/images.json", handleImages)

	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}
}
