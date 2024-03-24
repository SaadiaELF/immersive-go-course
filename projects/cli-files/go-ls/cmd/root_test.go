package cmd

import (
	"fmt"
	"os"
	"testing"
)

// go-ls tests
type testCase struct {
	path        []string
	expected    string
	description string
}

func Test_GetCurrentPath(t *testing.T) {
	expected, _ := os.Getwd()
	testCases := []testCase{
		{
			path:        []string{"cmd", "./assets"},
			expected:    "./assets",
			description: "Argument is provided",
		},
		{
			path:        []string{"cmd"},
			expected:    expected,
			description: "Argument is not provided",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d - %s", i+1, tc.description), func(t *testing.T) {
			os.Args = tc.path
			if path := getCurrentPath(); path != tc.expected {
				t.Errorf("Expected path: %s, got: %s", tc.expected, path)
			}
		})
	}
}
