package filtering

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// Create new filtering pipe
func TestFilteringPipe(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Test case 1", input: "Hello, World!", expected: "Hello, World!"},
		{name: "Test case 2", input: "Hello, 123 World!", expected: "Hello,  World!"},
		{name: "Test case 3", input: "Hello, 123 World! 456", expected: "Hello,  World! "},
		{name: "Test case 4", input: "Hello, 123 World! 456 789", expected: "Hello,  World!  "},
		{name: "Test case 5", input: "Hello, 123 World! 456=456", expected: "Hello,  World! ="},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var output bytes.Buffer
			filteringPipe := NewFilteringPipe(&output)
			filteringPipe.Write([]byte(testCase.input))

			require.Equal(t, testCase.expected, output.String())
		})
	}
}
