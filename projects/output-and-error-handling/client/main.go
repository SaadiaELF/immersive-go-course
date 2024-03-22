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
	client := &http.Client{}

	for retries := 0; retries <= maxRetries; retries++ {
		response, err := makeWeatherRequest(url, client, retries)
		if err != nil {
			handleError(err)
			os.Exit(1)
		}
		if response != "" {
			fmt.Fprintln(os.Stdout, response)
			os.Exit(0)
		}
	}
}

func makeWeatherRequest(url string, client *http.Client, retries int) (string, error) {
	resp, err := getWeather(url, client)
	if err != nil {
		return "", fmt.Errorf("failed to get weather: %w", err)
	}
	defer resp.Body.Close()

	return handleResponse(resp, retries)
}

func getWeather(url string, client *http.Client) (*http.Response, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make http request: %w", err)
	}
	return resp, nil
}

func handleResponse(resp *http.Response, retries int) (string, error) {
	switch resp.StatusCode {
	case 200:
		return handleSuccessResponse(resp)
	case 429:
		return "", handleRateLimited(resp.Header.Get("Retry-After"), retries)
	case 500:
		return "", fmt.Errorf("%d : Internal Server Error", resp.StatusCode)
	default:
		return "", fmt.Errorf("%d : Unexpected Error", resp.StatusCode)
	}
}
func handleSuccessResponse(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%d : Failed to read response body : %w", resp.StatusCode, err)
	}

	return string(body), nil
}
func handleError(err error) {
	fmt.Fprintln(os.Stderr, "Sorry we cannot get the weather!")
	fmt.Fprintf(os.Stderr, "%v\n", err)
}

func handleRateLimited(retryTime string, retries int) error {
	retrySeconds, err := convertTime(retryTime, RealTimeProvider{})
	if err != nil {
		return err
	}
	if retrySeconds > 1 && retrySeconds <= 5 && retries < maxRetries {
		fmt.Printf("We will retry to get you the weather. Please wait %d seconds\n", retrySeconds)
		time.Sleep(time.Duration(retrySeconds) * time.Second)
		return nil
	} else if retries >= maxRetries {
		return fmt.Errorf("internal Error : Failed to retry due to exceeded rate limit")
	} else if retrySeconds > 5 {
		return fmt.Errorf("internal Error : Failed to retry due to exceeded duration limit")
	} else {
		return fmt.Errorf("internal Error : Failed to retry  for an unknown reason")
	}
}

type TimeProvider interface {
	Now() time.Time
	Until(time.Time) time.Duration
}
type RealTimeProvider struct{}

func (r RealTimeProvider) Now() time.Time {
	return time.Now()
}

func (r RealTimeProvider) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func convertTime(retryTime string, tp TimeProvider) (int, error) {
	retrySeconds, err := strconv.Atoi(retryTime)
	if err == nil && retrySeconds != 0 {
		return retrySeconds, nil
	} else if _, err := time.Parse(http.TimeFormat, retryTime); err == nil {
		httpTime, err := time.Parse(http.TimeFormat, retryTime)
		if err != nil {
			return 0, fmt.Errorf("internal Error: error parsing HTTP Time Format")
		}
		retrySeconds := int(tp.Until(httpTime).Seconds())
		return retrySeconds, nil
	} else if retryTime == "a while" {
		return 5, nil
	} else {
		return 0, fmt.Errorf("internal Error: Failed to convert retry time")
	}
}
