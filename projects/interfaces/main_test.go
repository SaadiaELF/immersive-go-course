package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// If you make a buffer named b containing some bytes, calling b.Bytes() returns the same bytes you created it with.
func TestBufferBytesReturnsInitialBytes(t *testing.T) {
	inputBytes := []byte("Hello, World!")
	b := NewBuffer(inputBytes)

	bytes := b.Bytes()

	require.Equal(t, len(inputBytes), len(bytes))
}

// If you write some extra bytes to that buffer using b.Write(), a call to b.Bytes() returns both the initial bytes and the extra bytes.
func TestBufferWriteExtraBytes(t *testing.T) {
	inputBytes := []byte("Hello ")
	b := NewBuffer(inputBytes)
	n := b.Write([]byte("World!"))

	bytes := b.Bytes()

	require.Equal(t, len(inputBytes)+n, len(bytes))
}

// If you call b.Read() with a slice big enough to read all of the bytes in the buffer, all of the bytes are read.
func TestBufferReadFullSlice(t *testing.T) {
	inputBytes := []byte("Hello, World!")
	b := NewBuffer(inputBytes)

	outputBytes := make([]byte, len(inputBytes))
	n := b.Read(outputBytes)

	require.Equal(t, len(inputBytes), n)
}

// If you call b.Read() with a slice smaller than the contents of the buffer, some of the bytes are read. If you call it again, the next bytes are read.
func TestBufferReadPartialSlices(t *testing.T) {
	inputBytes := []byte("Hello, World!")
	b := bytes.NewBuffer(inputBytes)

	outputBytes := make([]byte, 7)
	b.Read(outputBytes)
	require.Equal(t, b.String(), "World!")
	b.Read(outputBytes)
	require.Equal(t, b.String(), "")
}
