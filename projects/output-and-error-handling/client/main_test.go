package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testCase struct {
	retryAfterTime string
	retries        int
	expectedValue  time.Duration
	expectedError  error
	description    string
}

func TestMakeWeatherRequest_Success(t *testing.T) {
	expRespSlice := [][]byte{[]byte("Today it will be sunny!"), []byte("Tomorrow it will be rainy!")}

	for i, expected := range expRespSlice {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(expected)
			}))
			defer svr.Close()
			client := &http.Client{}
			resp, err := makeWeatherRequest(client, svr.URL, 3)
			if err != nil {
				t.Errorf("expected err to be nil got %v", err)
			}

			if resp != string(expected) {
				t.Errorf("expected response body to be %s got %s", expected, resp)
			}
		})
	}
}

func TestMakeWeatherRequest_InternalServerError(t *testing.T) {
	errorMessage := "Internal Server Error"
	errorStatusCode := http.StatusInternalServerError
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(errorStatusCode)
		w.Write([]byte(errorMessage))
	}))
	defer svr.Close()

	client := &http.Client{}
	_, err := makeWeatherRequest(client, svr.URL, 3)
	if err == nil {
		t.Error("expected error, got nil")
	}

	expectedError := fmt.Sprint(errorStatusCode, " : ", errorMessage)
	if err.Error() != expectedError {
		t.Errorf("expected error message to be %s, got %s", expectedError, err.Error())
	}
}

func TestHandleRateLimited(t *testing.T) {
	testCases := []testCase{
		{
			retryAfterTime: "3",
			retries:        2,
			expectedValue:  3 * time.Second,
			expectedError:  nil,
			description:    "Success",
		},
		{
			retryAfterTime: "3",
			retries:        4,
			expectedValue:  3 * time.Second,
			expectedError:  fmt.Errorf("internal Error : Failed to retry due to exceeded rate limit"),
			description:    "Failure - Exceed rate limit",
		},
		{
			retryAfterTime: "6",
			retries:        2,
			expectedValue:  6 * time.Second,
			expectedError:  fmt.Errorf("internal Error : Failed to retry due to exceeded duration limit"),
			description:    "Failure - Exceed duration limit",
		},
		{
			retryAfterTime: "1",
			retries:        2,
			expectedValue:  1 * time.Second,
			expectedError:  fmt.Errorf("internal Error : Failed to retry  for an unknown reason"),
			description:    "Failure - Unknown reason",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d %s", i+1, tc.description), func(t *testing.T) {
			err := handleRateLimited(tc.retryAfterTime, tc.retries)
			if tc.expectedError == nil {
				if err != nil {
					t.Errorf("expected err to be nil got %v", err)
				}
			} else {
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error message to be %s, got %s", tc.expectedError, err.Error())
				}
			}
		})
	}
}

type MockTimeProvider struct{}

func (m MockTimeProvider) Now() time.Time {
	return time.Date(2024, time.March, 22, 12, 0, 0, 0, time.UTC)
}

func (m MockTimeProvider) Until(t time.Time) time.Duration {
	return t.Sub(m.Now())
}

func TestConvertTime(t *testing.T) {
	tp := MockTimeProvider{}

	testCases := []testCase{
		{
			retryAfterTime: "3",
			expectedValue:  3 * time.Second,
			expectedError:  nil,
			description:    "Success - Integer numbers of seconds",
		},
		{
			retryAfterTime: (tp.Now().Add(3 * time.Second)).Format(http.TimeFormat),
			expectedValue:  3 * time.Second,
			expectedError:  nil,
			description:    "Success - Timestamps",
		},
		{
			retryAfterTime: "a while",
			expectedValue:  5 * time.Second,
			expectedError:  nil,
			description:    "Success - A while",
		},
		{
			retryAfterTime: "Invalid Format",
			expectedValue:  0 * time.Second,
			expectedError:  fmt.Errorf("internal Error: Failed to convert retry time"),
			description:    "Failure - Invalid Format",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d %s", i+1, tc.description), func(t *testing.T) {

			seconds, err := convertTime(tp, tc.retryAfterTime)
			if tc.expectedError == nil {
				if seconds != tc.expectedValue {
					t.Errorf("expected seconds to equal to %d got %d", seconds, tc.expectedValue)
				}
				if err != nil {
					t.Errorf("expected err to be nil got %v", err)
				}
			} else {
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("expected error message to be %s, got %s", tc.expectedError, err.Error())
				}
			}

		})
	}
}
