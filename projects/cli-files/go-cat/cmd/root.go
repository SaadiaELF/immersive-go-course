package cmd

import "fmt"

func Execute() {
	args := checkArgs([]string{})
	fmt.Println(args)
}

func checkArgs(args []string) (err error) {
	if len(args) == 0 {
		err = fmt.Errorf("error: no file specified")
	}
	return err
}
