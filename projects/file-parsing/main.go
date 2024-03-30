package main

import (
	"fmt"
	"os"

	jsonparser "github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/fileparser"
)

func main() {
	data, err := jsonparser.ParseJSON("./examples/json.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(data)
}
