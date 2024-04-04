package binary

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestByteOrder(t *testing.T) {
	testCases := []struct {
		filename string
		expected binary.ByteOrder
	}{
		{filename: "../../examples/custom-binary-le.bin", expected: binary.LittleEndian},
		{filename: "../../examples/custom-binary-be.bin", expected: binary.BigEndian},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			file, _ := os.Open(tc.filename)
			binaryOrder, err := ByteOrder(file)
			require.NoError(t, err)
			require.Equal(t, tc.expected, binaryOrder)
		})
	}
}

func TestBinaryParser(t *testing.T) {
	players := Players{{Name: "Aya", HighScore: 10}, {Name: "Prisha", HighScore: 30}, {Name: "Charlie", HighScore: -1}, {Name: "Margot", HighScore: 25}}
	testCases := []struct {
		filename string
		expected Players
	}{
		{filename: "../../examples/custom-binary-le.bin", expected: players},
		{filename: "../../examples/custom-binary-be.bin", expected: players},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			players, err := BinaryParser(tc.filename)
			require.NoError(t, err)
			require.Equal(t, tc.expected, players)
		})
	}
}
