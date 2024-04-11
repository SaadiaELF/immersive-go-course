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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/images.json", handleImages)

	fmt.Fprintln(os.Stderr, "Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen: %v", err)
		os.Exit(1)
	}
}

func handleImages(w http.ResponseWriter, r *http.Request) {
	conn, err := setDatabaseConnection()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	defer conn.Close(context.Background())
	images, err := fetchImages(conn)
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

func fetchImages(conn *pgx.Conn) ([]types.Image, error) {
	var images []types.Image
	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")
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
