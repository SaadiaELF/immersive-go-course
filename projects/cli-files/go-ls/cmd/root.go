package cmd

import (
	"flag"
	"fmt"
	"os"
)

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

func Execute() {
	var path string

	h := flag.Bool("h", false, "show description and usage")
	flag.Parse()

	if *h {
		flag.Usage()
		return
	}

	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path, _ = os.Getwd()
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fdopendir := fmt.Sprintf("fdopendir %s: not a directory", path)
		if err.Error() == fdopendir {
			fmt.Println(path)
		} else {
			fmt.Println(err)
		}
		return
	}
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}

// }
