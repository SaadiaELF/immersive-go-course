package static

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(path *string, port *string) {

	http.Handle("/", http.FileServer(http.Dir(*path)))

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
