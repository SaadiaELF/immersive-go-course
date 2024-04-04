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
