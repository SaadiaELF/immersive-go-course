package binary

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/types"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/utils"
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
	myParser := Parser{}
	testCases := []struct {
		name     string
		filename string
		expected types.Players
	}{
		{name: "LittleEndian", filename: "../../examples/custom-binary-le.bin", expected: utils.Players},
		{name: "BigEndian", filename: "../../examples/custom-binary-be.bin", expected: utils.Players},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			players, err := myParser.Parse(tc.filename)
			require.NoError(t, err)
			require.Equal(t, tc.expected, players)
			require.Equal(t, tc.expected, players)
		})
	}
}
