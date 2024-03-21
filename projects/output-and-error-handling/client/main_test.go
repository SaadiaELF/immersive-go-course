package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
			resp, err := handleWeatherRequest(svr.URL, client, time.Now(), 3)
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
	_, err := handleWeatherRequest(svr.URL, client, time.Now(), 3)
	if err == nil {
		t.Error("expected error, got nil")
	}

	expectedError := fmt.Sprint(errorStatusCode, " : ", errorMessage)
	if err.Error() != expectedError {
		t.Errorf("expected error message to be %s, got %s", expectedError, err.Error())
	}
}

func TestHandleRateLimited(t *testing.T) {
	testCases := []int{3, 6, -1}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			retryAfterHeader := strconv.Itoa(tc)
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Retry-After", retryAfterHeader)
			}))
			defer svr.Close()

			err := handleRateLimited(retryAfterHeader, time.Now(), 2)
			if tc == 3 {
				if err != nil {
					t.Errorf("expected err to be nil got %v", err)
				}
			}
			if tc == 6 || tc == -1 {
				if err == nil {
					t.Error("expected error, got nil")
				}

				expectedError := "internal Error : Failed to retry"
				if err.Error() != expectedError {
					t.Errorf("expected error message to be %s, got %s", expectedError, err.Error())
				}
			}
		})
	}
}
