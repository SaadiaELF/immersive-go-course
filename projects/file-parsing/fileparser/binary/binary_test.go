package binary

import (
	"encoding/binary"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestByteOrder(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		expected binary.ByteOrder
	}{
		{name: "LittleEndian", filename: "../../examples/custom-binary-le.bin", expected: binary.LittleEndian},
		{name: "BigEndian", filename: "../../examples/custom-binary-be.bin", expected: binary.BigEndian},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, _ := os.Open(tc.filename)
			binaryOrder, err := ByteOrder(file)
			require.NoError(t, err)
			require.Equal(t, tc.expected, binaryOrder)
		})
	}
}

func TestBinaryParser(t *testing.T) {
	testCases := []struct {
		name                     string
		filename                 string
		expected_HighScorePlayer string
		expected_LowScorePlayer  string
	}{
		{name: "LittleEndian", filename: "../../examples/custom-binary-le.bin", expected_HighScorePlayer: "Prisha", expected_LowScorePlayer: "Charlie"},
		{name: "BigEndian", filename: "../../examples/custom-binary-be.bin", expected_HighScorePlayer: "Prisha", expected_LowScorePlayer: "Charlie"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			highestScorePlayer, lowestScorePlayer, err := BinaryParser(tc.filename)
			require.ErrorIs(t, err, io.EOF)
			require.Equal(t, tc.expected_HighScorePlayer, highestScorePlayer)
			require.Equal(t, tc.expected_LowScorePlayer, lowestScorePlayer)
		})
	}
}
