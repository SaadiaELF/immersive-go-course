package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Hello...")

	// Create a channel to receive signals
	sigCh := make(chan os.Signal, 1)

	// Notify the channel to receive SIGINT signals
	signal.Notify(sigCh, syscall.SIGINT)

	// Wait for the signal
	<-sigCh

	// When the signal is received, print the message
	fmt.Println("\rCtrl+C is not responding")

	// Wait for the signal again
	<-sigCh
}
