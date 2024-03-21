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
	currentTime := time.Now()
	// Create an HTTP client
	client := &http.Client{}
	// This loop will run only if the number of retries doesn't exceed 3
	for retries := 0; retries <= maxRetries; retries++ {
		resp, err := getWeather(url, client)
		if err != nil {
			handleError(err)
			os.Exit(1)
		}
		// Close the response body after been fully read
		defer resp.Body.Close()

		if err := handleResponse(resp, currentTime, retries); err != nil {
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
func handleResponse(resp *http.Response, currentTime time.Time, retries int) error {
	switch resp.StatusCode {
	case 200:
		return handleSuccessResponse(resp)
	case 429:
		return handleRateLimited(resp.Header.Get("Retry-After"), currentTime, retries)
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
func handleRateLimited(retryTime string, currentTime time.Time, retries int) error {
	retrySeconds, err := convertTime(retryTime, currentTime)
	if err != nil {
		return err
	}
	if retrySeconds > 1 && retrySeconds <= 5 && retries < maxRetries {
		fmt.Printf("We will retry to get you the weather. Please wait %d seconds\n", retrySeconds)
		time.Sleep(time.Duration(retrySeconds) * time.Second)
	} else {
		return fmt.Errorf("internal Error : Failed to retry")
	}
	return nil
}

func convertTime(retryTime string, currentTime time.Time) (int, error) {

	if retryTime == http.TimeFormat {
		httpTime, err := time.Parse(http.TimeFormat, retryTime)
		if err != nil {
			return 0, fmt.Errorf("error parsing HTTP Time Format: %v", err)
		}
		return int(currentTime.Sub(httpTime).Seconds()), nil
	}

	if retryTime == "a while" {
		return 5, nil
	} else {
		retrySeconds, err := strconv.Atoi(retryTime)
		if err != nil {
			return 0, fmt.Errorf("internal Error: Failed to convert retry time: %v", err)
		}
		return retrySeconds, nil
	}
}
