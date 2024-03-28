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
		err := checkArgs(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fileLines, err := readFileLines(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, fileLine := range fileLines {
			fmt.Print(fileLine)
		}
	}
}

func checkArgs(arg string) error {
	if arg == "" {
		return fmt.Errorf("go-cat: no file specified")
	}
	fileInfo, err := os.Stat(arg)
	if err != nil {
		return fmt.Errorf("go-cat: '%s': no such file or directory", arg)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("go-cat: '%s': is a directory", arg)
	}
	return nil
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
			fileLines = append(fileLines, "\n")
			break
		}
		fileLines = append(fileLines, string(buffer[0:n]))
	}
	return fileLines, nil
}
