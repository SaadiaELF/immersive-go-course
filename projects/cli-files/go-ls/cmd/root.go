package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	// Open the current directory
	dir, err := os.Open(".")
	if err != nil {
		fmt.Printf("Error opening directory: %v\n", err)
		return
	}
	defer dir.Close()

	// Read directory all entries
	entries, err := dir.ReadDir(0)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}
