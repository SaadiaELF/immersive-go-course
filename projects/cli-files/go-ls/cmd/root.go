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
		fdopendir := fmt.Sprintf("fdopendir %s: not a directory", path)
		if err.Error() == fdopendir {
			fmt.Println(path)
		} else {
			fmt.Println(err)
		}
		return
	}
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}
