package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	args := os.Args[1:]
	checkArgs(args)
	fmt.Println(checkArgs(args))
}

func checkArgs(args []string) (err error) {

	if len(args) == 1 {
		fileInfo, _ := os.Stat(args[0])
		if fileInfo.IsDir() {
			err = fmt.Errorf("error: '%s': is a directory", args[0])
		}
	}
	if len(args) == 0 {
		err = fmt.Errorf("error: no file specified")
	}
	return err
}
