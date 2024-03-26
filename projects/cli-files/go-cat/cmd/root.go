package cmd

import (
	"fmt"
	"os"
)

func Execute() {
	path, err := checkArg()
	if err != nil {
		fmt.Println(err)
	} else {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		os.Stdout.Write(data)
	}

}

func checkArg() (string, error) {
	if len(os.Args) > 1 {
		return os.Args[1], nil
	} else {
		err := os.ErrInvalid
		return "", err
	}
}
