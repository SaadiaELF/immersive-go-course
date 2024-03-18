package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	retries := 0
	for retries <= 3 {
		// Connect to the server and getting a response
		resp, err := http.Get("http://localhost:8080")
		// Show error message if connection is not established
		if err != nil {
			fmt.Print("Sorry we cannot get the weather!\n")
			fmt.Fprintf(os.Stderr, "Failed to make http request: %v\n", err)
			os.Exit(1)
		}

		// Close the response body after been fully read
		defer resp.Body.Close()

		// Read response body and store it in a var
		body, err := io.ReadAll(resp.Body)

		// Show error message if we cannot read the response body
		if err != nil {
			fmt.Print("Sorry we cannot get the weather!\n")
			fmt.Fprintf(os.Stderr, "Failed to read response body: %v\n", err)
			os.Exit(1)
		}

		// Handle cases depending on the Status code of the response
		switch resp.StatusCode {
		case 200:
			// Convert response body from binary to string
			retries = 0
			sb := string(body)
			fmt.Fprint(os.Stdout, sb+"\n")
			os.Exit(0)
		case 429:
			handleRateLimited(resp.Header.Get("Retry-After"), retries)
		case 500:
			fmt.Print("Sorry we cannot get the weather!\n")
			fmt.Fprint(os.Stderr, "Internal Server Error\n")
			os.Exit(1)
		default:
			fmt.Print("Sorry we cannot get the weather!\n")
			fmt.Fprint(os.Stderr, "Unexpected Error\n")
			os.Exit(1)
		}
	}
}

// Handle response and retry depending on the Retry-After header
func handleRateLimited(retryTime string, retries int) {
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
			fmt.Print("Sorry we cannot get the weather!")
			fmt.Fprintf(os.Stderr, "Internal Error : Failed to convert retry time: %v\n", err)
			os.Exit(1)
		}
	}
	if retrySeconds > 1 && retrySeconds <= 5 {
		fmt.Printf("We will retry to get you the weather. Please wait %d seconds\n", retrySeconds)
		time.Sleep(time.Duration(retrySeconds) * time.Second)
		retries++
	} else {
		fmt.Print("Sorry we cannot get the weather!\n")
		fmt.Fprint(os.Stderr, "Internal Error : Failed to retry\n")
		os.Exit(1)
	}
}
