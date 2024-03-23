package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	// Read directory all entries
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}
