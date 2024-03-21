package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	url        = "http://localhost:8080"
	maxRetries = 3
)

func main() {

	// Create an HTTP client
	client := &http.Client{}
	// This loop will run only if the number of reties doesn't exceed 3
	for retries := 0; retries < maxRetries; retries++ {
		resp, err := getWeather(url, client)
		if err != nil {
			handleError(err)
			os.Exit(1)
		}
		// Close the response body after been fully read
		defer resp.Body.Close()

		if err := handleResponse(resp); err != nil {
			handleError(err)
			os.Exit(1)
		}
	}
}

func getWeather(url string, client *http.Client) (*http.Response, error) {
	// Connect to the server and getting a response
	resp, err := client.Get(url)
	// Show error message if connection is not established
	if err != nil {
		return nil, fmt.Errorf("failed to make http request: %v", err)
	}
	return resp, nil
}

// Handle cases depending on the Status code of the response
func handleResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case 200:
		return handleSuccessResponse(resp)
	case 429:
		return handleRateLimited(resp.Header.Get("Retry-After"))
	case 500:
		return fmt.Errorf("%d : Internal Server Error", resp.StatusCode)
	default:
		return fmt.Errorf("%d : Unexpected Error", resp.StatusCode)
	}
}
func handleSuccessResponse(resp *http.Response) error {
	// Read response body and store it in a var
	body, err := io.ReadAll(resp.Body)
	// Show error message if we cannot read the response body
	if err != nil {
		return fmt.Errorf("%d : Failed to read response body : %v", resp.StatusCode, err)
	}
	// Convert response body from binary to string
	fmt.Fprintln(os.Stdout, string(body))
	return nil
}
func handleError(err error) {
	fmt.Fprintln(os.Stderr, "Sorry we cannot get the weather!")
	fmt.Fprintf(os.Stderr, "%v\n", err)
}

// Handle response and retry depending on the Retry-After header
func handleRateLimited(retryTime string) error {
	retrySeconds := 0
	var err error
	retryTimeDate, err := time.Parse(time.RFC1123, retryTime)

	if err == nil {
		retrySeconds = int(time.Until(retryTimeDate).Seconds())
	} else if retryTime == "a while" {
		retrySeconds = 5
	} else {
		retrySeconds, err = strconv.Atoi(retryTime)
		if err != nil {
			return fmt.Errorf("internal Error: Failed to convert retry time: %v", err)
		}
	}
	if retrySeconds > 1 && retrySeconds <= 5 {
		fmt.Printf("We will retry to get you the weather. Please wait %d seconds\n", retrySeconds)
		time.Sleep(time.Duration(retrySeconds) * time.Second)
	} else {
		return fmt.Errorf("internal Error : Failed to retry")
	}
	return nil
}
