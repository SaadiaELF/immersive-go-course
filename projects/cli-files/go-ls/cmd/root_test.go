package cmd

import (
	"fmt"
	"os"
	"testing"
)

// go-ls tests

func Test_GetCurrentPath(t *testing.T) {
	expectedPath, _ := os.Getwd()
	testCases := []struct {
		path        []string
		expected    string
		description string
	}{
		{
			path:        []string{"cmd", "./assets"},
			expected:    "./assets",
			description: "Argument is provided",
		},
		{
			path:        []string{"cmd"},
			expected:    expectedPath,
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

func Test_Listcontent(t *testing.T) {

}
