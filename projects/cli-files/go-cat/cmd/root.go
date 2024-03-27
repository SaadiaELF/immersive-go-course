package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	args := os.Args[1:]
	err := checkArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func checkArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("error: no file specified")
	}
	if len(args) == 1 {
		fileInfo, err := os.Stat(args[0])
		if err != nil {
			return fmt.Errorf("error: '%s': no such file or directory", args[0])
		}
		if fileInfo.IsDir() {
			return fmt.Errorf("error: '%s': is a directory", args[0])
		}
	}
	return nil
}
