package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Hello...")

	var pid int
	// Create a channel to receive signals
	sigCh := make(chan os.Signal, 1)

	// Notify the channel to receive SIGINT signals
	signal.Notify(sigCh, syscall.SIGINT)

	// Wait for the signal
	<-sigCh

	// When the signal is received, print the message
	log.Println("\rEnter the PID of the process you want to kill: ")

	// Scan the standard input for the PID
	_, err := fmt.Scan(&pid)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Kill the process
	err = syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	// Print success message
	log.Printf("Process %v killed\n", pid)
}
