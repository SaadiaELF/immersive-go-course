package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func Execute() {
	args := os.Args[1:]
	err := checkArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fileLines, err := readFile(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	for _, fileLine := range fileLines {
		fmt.Print(fileLine)
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

func readFile(filePath string) ([]string, error) {
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
			break
		}
		fileLines = append(fileLines, string(buffer[0:n]))
	}
	return fileLines, nil
}
