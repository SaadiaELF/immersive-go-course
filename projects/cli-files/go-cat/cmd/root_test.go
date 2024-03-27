package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// cat [file-path] displays the content of a file
// cat [directory-path] displays error message : [directory-path] is a directory
// cat [non-existent-file-path] displays error message : [non-existent-file-path]: no such file or directory
// cat [file-path] [file-path] displays the content of both files
// cat no ARG displays error message : no file specified
// Extra test cases:
// cat [file-path] [directory-path] displays error message : [directory-path] is a directory
// cat [file-path] [non-existent-file-path] displays error message : [non-existent-file-path]: no such file or directory
// cat [directory-path] [file-path] displays error message : [directory-path] is a directory

func TestNoFileSpecified(t *testing.T) {
	// TestNoFileSpecified tests the case when no file is specified
	// Expected output: "no file specified"
	err := checkArgs([]string{})
	require.Equal(t, "error: no file specified", err)
}
