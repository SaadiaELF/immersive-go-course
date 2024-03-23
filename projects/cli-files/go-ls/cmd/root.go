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
	
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(path)
		return
	}
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}
