package e2e_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestRoutes(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("err loading: %v", err)
	}
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "80"
	}
	testcases := []struct {
		name     string
		endpoint string
		expected string
	}{
		{
			name:     "root",
			endpoint: "/",
			expected: "Hello, world!\n",
		},
		{
			name:     "ping",
			endpoint: "/ping",
			expected: "pong pong ...\n",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s%s", port, tc.endpoint), nil)
			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatalf("Could not make request: %v", err)
			}

			if response.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", response.StatusCode)
			}
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("Could not read response body: %v", err)
			}
			if string(body) != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, string(body))
			}

			defer response.Body.Close()
		})
	}

}
