package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// cat [file-path] displays the content of a file
// cat [directory-path] displays error message : [directory-path] is a directory
// cat [file-path] [file-path] displays the content of both files
// cat no ARG displays error message : no file specified
// Extra test cases:
// cat [file-path] [directory-path] displays error message : [directory-path] is a directory
// cat [file-path] [non-existent-file-path] displays error message : [non-existent-file-path]: no such file or directory
// cat [directory-path] [file-path] displays error message : [directory-path] is a directory

// cat [non-existent-file-path] displays error message : [non-existent-file-path]: no such file or directory
func TestNoFileSpecified(t *testing.T) {
	err := checkArgs([]string{})
	require.Equal(t, err.Error(), "error: no file specified")
}
