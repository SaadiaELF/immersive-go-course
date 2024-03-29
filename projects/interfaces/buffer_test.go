package OurBuffer

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// If you make a buffer named b containing some bytes, calling b.Bytes() returns the same bytes you created it with.
func TestBufferBytesReturnsInitialBytes(t *testing.T) {
	inputBytes := []byte("Hello, World!")
	b := NewBuffer(inputBytes)

	bytes := b.Bytes()

	excepted := string(inputBytes)
	actual := string(bytes)

	require.Equal(t, excepted, actual)
}

// If you write some extra bytes to that buffer using b.Write(), a call to b.Bytes() returns both the initial bytes and the extra bytes.
func TestBufferWriteExtraBytes(t *testing.T) {
	inputBytes := []byte("Hello ")
	b := NewBuffer(inputBytes)
	n, err := b.Write([]byte("World!"))
	require.NoError(t, err)

	bytes := b.Bytes()
	expected := string((bytes))
	actual := string(inputBytes) + "World!"
	require.Equal(t, expected, actual)
	require.Equal(t, len(bytes), len(inputBytes)+n)
}

// If you call b.Read() with a slice big enough to read all of the bytes in the buffer, all of the bytes are read.
func TestBufferReadFullSlice(t *testing.T) {
	inputBytes := []byte("Hello, World!")
	b := NewBuffer(inputBytes)

	outputBytes := make([]byte, len(inputBytes))
	n, err := b.Read(outputBytes)
	require.ErrorAs(t, err, &io.EOF)

	expected := string(inputBytes)
	actual := string(outputBytes)
	require.Equal(t, expected, actual)
	require.Equal(t, len(inputBytes), n)
}

// If you call b.Read() with a slice smaller than the contents of the buffer, some of the bytes are read. If you call it again, the next bytes are read.
func TestBufferReadPartialSlices(t *testing.T) {
	inputBytes := []byte("Hello, World!")
	b := NewBuffer(inputBytes)

	outputBytes := make([]byte, 7)
	t.Run("first read", func(t *testing.T) {
		n, err := b.Read(outputBytes)
		require.NoError(t, err)
		require.Equal(t, "World!", b.String())
		require.Equal(t, "Hello, ", string(outputBytes))
		require.Equal(t, 7, n)
	})

	t.Run("Second Read", func(t *testing.T) {
		n, err := b.Read(outputBytes)
		require.ErrorAs(t, err, &io.EOF)
		require.Equal(t, b.String(), "")
		require.Equal(t, "World! ", string(outputBytes))
		require.Equal(t, 6, n)
	})
}
