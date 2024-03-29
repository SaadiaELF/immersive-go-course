package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// If you write some extra bytes to that buffer using b.Write(), a call to b.Bytes() returns both the initial bytes and the extra bytes.
// If you call b.Read() with a slice big enough to read all of the bytes in the buffer, all of the bytes are read.
// If you call b.Read() with a slice smaller than the contents of the buffer, some of the bytes are read. If you call it again, the next bytes are read.

// If you make a buffer named b containing some bytes, calling b.Bytes() returns the same bytes you created it with.
func TestBytesMethod(t *testing.T) {
	// Create a buffer named b containing some bytes.
	var b bytes.Buffer
	n, _ := b.Write([]byte("Hello, World!"))
	// Call b.Bytes() and check that it returns the same bytes you created it with.
	bytes := b.Bytes()

	require.Equal(t, n, len(bytes))

}
