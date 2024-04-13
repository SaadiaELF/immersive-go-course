package static

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	// Parse command line flags
	path := flag.String("path", "./", "path where the static files are read from")
	port := flag.String("port", "8080", "port to run the static server")
	flag.Parse()

	//build the server and start it
	handler := http.FileServer(http.Dir(*path))
	http.Handle("/", handler)

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
