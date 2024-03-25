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
		response, err := makeWeatherRequest(client, url, retries)
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

func makeWeatherRequest(client *http.Client, url string, retries int) (string, error) {
	resp, err := getWeather(client, url)
	if err != nil {
		return "", fmt.Errorf("failed to get weather: %w", err)
	}
	defer resp.Body.Close()

	return handleResponse(resp, retries)
}

func getWeather(client *http.Client, url string) (*http.Response, error) {
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
	retryDuration, err := convertTime(RealTimeProvider{}, retryTime)
	if err != nil {
		return err
	}
	switch {
	case retryDuration > 1*time.Second && retryDuration <= 5*time.Second && retries < maxRetries:
		fmt.Printf("We will retry to get you the weather. Please wait %v.\n", retryDuration)
		time.Sleep(retryDuration)
		return nil
	case retries >= maxRetries:
		return fmt.Errorf("internal Error : Failed to retry due to exceeded rate limit")
	case retryDuration > 5*time.Second:
		return fmt.Errorf("internal Error : Failed to retry due to exceeded duration limit")
	default:
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

func convertTime(tp TimeProvider, retryTime string) (time.Duration, error) {
	secondsInt, err := strconv.Atoi(retryTime)
	if err == nil && secondsInt != 0 {
		retryDuration := time.Duration(secondsInt) * time.Second
		return retryDuration, nil
	} else if _, err := time.Parse(http.TimeFormat, retryTime); err == nil {
		httpTime, err := time.Parse(http.TimeFormat, retryTime)
		if err != nil {
			return 0, fmt.Errorf("internal Error: error parsing HTTP Time Format")
		}
		retryDuration := tp.Until(httpTime)
		return retryDuration, nil
	} else if retryTime == "a while" {
		return 5 * time.Second, nil
	} else {
		return 0, fmt.Errorf("internal Error: Failed to convert retry time")
	}
}
