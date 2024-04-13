package static

import (
	"flag"
	"fmt"
)

func Run() {
	path := flag.String("path", "./", "path where the static files are read from")
	port := flag.String("port", "8080", "port to run the static server")
	flag.Parse()
	fmt.Println("path:", *path)
	fmt.Println("port:", *port)
	fmt.Println("Hello from static server!")
}
