package main

import (
	"flag"
	"servers/static"
)

func main() {
	// Parse command line flags
	path := flag.String("path", "./", "path where the static files are read from")
	port := flag.String("port", "8080", "port to run the static server")
	flag.Parse()
	static.Run(path, port)
}
