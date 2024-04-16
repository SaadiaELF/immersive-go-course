package main

import (
	"flag"
	"log"
	"os"
	"servers/api"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	databaseURL := os.Getenv("DATABASE_URL")

	// Parse command line flags
	dbURL := flag.String("path", databaseURL, "path where the static files are read from")
	port := flag.String("port", "8081", "port to run the static server")
	flag.Parse()
	api.Run(dbURL, port)
}
