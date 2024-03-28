package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func Execute() {
	args := os.Args[1:]

	for _, arg := range args {
		fileLines, err := readFileLines(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "go-cat:", err)
			os.Exit(1)
		}
		for _, fileLine := range fileLines {
			fmt.Print(fileLine)
		}
	}
}

func readFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 16)
	fileLines := []string{}

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			if n == 0 {
				break
			}
		}
		fileLines = append(fileLines, string(buffer[0:n]))
	}
	return fileLines, nil
}
