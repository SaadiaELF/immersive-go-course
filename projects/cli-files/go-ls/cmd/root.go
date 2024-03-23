package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	} else {
		arg, _ = os.Getwd()
	}
	if arg == "-h" {
		fmt.Println("go-ls : list directory contents")
	} else {
		entries, err := os.ReadDir(arg)
		if err != nil {
			fdopendir := fmt.Sprintf("fdopendir %s: not a directory", arg)
			if err.Error() == fdopendir {
				fmt.Println(arg)
			} else {
				fmt.Println(err)
			}
			return
		}
		for _, entry := range entries {
			fmt.Println(entry.Name())
		}
	}
}
