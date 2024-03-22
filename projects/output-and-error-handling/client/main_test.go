package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleWeatherRequest_Success(t *testing.T) {
	expRespSlice := [][]byte{[]byte("Today it will be sunny!"), []byte("Tomorrow it will be rainy!")}

	for i, expected := range expRespSlice {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(expected)
			}))
			defer svr.Close()
			client := &http.Client{}
			resp, err := makeWeatherRequest(svr.URL, client, 3)
			if err != nil {
				t.Errorf("expected err to be nil got %v", err)
			}

			if resp != string(expected) {
				t.Errorf("expected response body to be %s got %s", expected, resp)
			}
		})
	}
}

func TestHandleWeatherRequest_InternalServerError(t *testing.T) {
	errorMessage := "Internal Server Error"
	errorStatusCode := http.StatusInternalServerError
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(errorStatusCode)
		w.Write([]byte(errorMessage))
	}))
	defer svr.Close()

	client := &http.Client{}
	_, err := makeWeatherRequest(svr.URL, client, 3)
	if err == nil {
		t.Error("expected error, got nil")
	}

	expectedError := fmt.Sprint(errorStatusCode, " : ", errorMessage)
	if err.Error() != expectedError {
		t.Errorf("expected error message to be %s, got %s", expectedError, err.Error())
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
	type testCase struct {
		retryAfterTime string
		expectedValue  int
		expectedError  error
	}

	testCases := []testCase{
		{
			retryAfterTime: "3",
			expectedValue:  3,
			expectedError:  nil,
		},
		{
			retryAfterTime: "6",
			expectedValue:  6,
			expectedError:  nil,
		},
		{
			retryAfterTime: (tp.Now().Add(3 * time.Second)).Format(http.TimeFormat),
			expectedValue:  3,
			expectedError:  nil,
		},
		{
			retryAfterTime: (tp.Now().Add(6 * time.Second)).Format(http.TimeFormat),
			expectedValue:  6,
			expectedError:  nil,
		},
		{
			retryAfterTime: "a while",
			expectedValue:  5,
			expectedError:  nil,
		},
		{
			retryAfterTime: "Invalid Format",
			expectedValue:  0,
			expectedError:  fmt.Errorf("internal Error: Failed to convert retry time"),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {

			seconds, err := convertTime(tc.retryAfterTime, tp)
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
