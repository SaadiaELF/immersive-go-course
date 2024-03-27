package cmd

import "fmt"

func Execute() {
	args := checkArgs([]string{})
	fmt.Println(args)
}

func checkArgs(args []string) (output string) {
	if len(args) == 0 {
		output = "error: no file specified"
	}
	return output
}
