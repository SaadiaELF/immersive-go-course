package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"
)

func Execute() {
	defer addFlag()
	listContent()
}

func boldText(s string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", s)
}

func init() {
	flag.Usage = func() {
		fmt.Printf("%s\n%s\n	go-ls - list directory content\n%s\n	go-ls [path] || go-ls [options]\n", boldText("Usage of go-ls"), boldText("Name:"), boldText("Format:"))
		fmt.Printf("%s\n	For each operand that names a file of a type other than directory, go-ls displays its name.\n	For each operand that names a file of type directory, go-ls displays the names of files contained within that directory.\n	If no operands are given, the contents of the current directory are displayed.\n", boldText("Description:"))
		fmt.Println("The following options are available:")
		flag.PrintDefaults()
	}
}
func addFlag() {
	h := flag.Bool("h", false, "show description and usage")
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
}
func listContent() {
	path := getCurrentPath()
	entries, err := os.ReadDir(path)
	if err != nil {
		handleErrors(err, path)
	}
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}
func getCurrentPath() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	} else {
		path, _ := os.Getwd()
		return path
	}
}
func handleErrors(err error, path string) {
	if errors.Is(err, syscall.ENOTDIR) {
		fmt.Println(path)
	} else {
		fmt.Println(err)
	}
}
