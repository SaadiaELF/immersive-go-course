package images

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

// POST /images.json
func PostImage(dbPool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {

	// Parse the request body
	var image Image
	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Error: failed to decode image: %v\n", err)
	}

	// Check if the image url exists in the database
	query := "SELECT url FROM public.images WHERE url = $1"
	row := (*dbPool).QueryRow(context.Background(), query, image.URL)
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
func FetchImages(pool *pgxpool.Pool, limit int) ([]Image, error) {
	var images []Image
	query := fmt.Sprintf("SELECT title, url, alt_text FROM public.images LIMIT %d", limit)
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to fetch images: %v\n", err)
	}
	for rows.Next() {
		var title, url, altText string
		err = rows.Scan(&title, &url, &altText)
		images = append(images, Image{Title: title, URL: url, AltText: altText})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to scan image: %v\n", err)
		}
	}
	return images, err
}

// getIndentParam returns the indent parameter from the request query string.
func GetIndentParam(r *http.Request) (string, error) {
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
