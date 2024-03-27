package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// cat [file-path] [file-path] displays the content of both files
// Extra test cases:
// cat [file-path] [directory-path] displays error message : [directory-path] is a directory
// cat [file-path] [non-existent-file-path] displays error message : [non-existent-file-path]: no such file or directory
// cat [directory-path] [file-path] displays error message : [directory-path] is a directory

// cat [non-existent-file-path] displays error message : [non-existent-file-path]: no such file or directory
func TestNoFileSpecified(t *testing.T) {
	err := checkArgs([]string{})
	require.Equal(t, err.Error(), "error: no file specified")
}

// cat [directory-path] displays error message : [directory-path] is a directory
func TestDirectoryPath(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{{name: "current directory", args: []string{"."}}, {name: "directory path", args: []string{"../../assets"}}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkArgs(tc.args)
			require.Equal(t, "error: '"+tc.args[0]+"': is a directory", err.Error())
		})
	}
}

// cat [non-existent-file-path] displays error message : [non-existent-file-path]: no such file or directory
func TestNonExistentFilePath(t *testing.T) {
	err := checkArgs([]string{"non-existent-file-path"})
	require.Equal(t, "error: 'non-existent-file-path': no such file or directory", err.Error())
}

func createTempFile() string {
	file, err := os.CreateTemp(".", "temp.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	_, err = file.WriteString("file contents") // less than 16 bytes
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return file.Name()
}

// cat [file-path] displays the content of a file
func TestValidFilePath(t *testing.T) {
	fileName := createTempFile()
	fileLines, err := readFile(fileName)
	require.NoError(t, err)
	require.Equal(t, []string{"file contents"}, fileLines)
	defer os.Remove(fileName)
}
