package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// cat [directory-path] displays error message : [directory-path] is a directory
func TestDirectoryPath(t *testing.T) {
	testCases := []struct {
		name string
		arg  string
	}{{name: "current directory", arg: "."}, {name: "directory path", arg: "../../assets"}}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := readFileLines(tc.arg)
			require.Equal(t, "read "+tc.arg+": is a directory", err.Error())
		})
	}
}

// cat [non-existent-file-path] displays error message : [non-existent-file-path]: no such file or directory
func TestNonExistentFilePath(t *testing.T) {
	fileName, err := createTempFile()
	require.NoError(t, err)
	os.Remove(fileName)
	_, err = readFileLines(fileName)
	require.Equal(t, "open "+fileName+": no such file or directory", err.Error())
}

func createTempFile() (string, error) {
	file, err := os.CreateTemp(".", "temp.txt")
	if err != nil {
		return "", err
	}
	_, err = file.WriteString("file contents") // less than 16 bytes
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

// cat [file-path] displays the content of a file
func TestValidFilePath(t *testing.T) {
	fileName, err := createTempFile()
	require.NoError(t, err)
	defer os.Remove(fileName)

	fileLines, err := readFileLines(fileName)
	require.NoError(t, err)
	require.Equal(t, []string{"file contents"}, fileLines)
}
