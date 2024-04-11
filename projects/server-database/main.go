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

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	// Create instance of the server
	server := &http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/images.json", handleImages)

	go func() {
		fmt.Fprintln(os.Stderr, "Listening on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: failed to listen and serve: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

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
