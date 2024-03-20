package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeather_Success(t *testing.T) {
	expRespSlice := [][]byte{[]byte("Today it will be sunny!"), []byte("Tomorrow it will be rainy!")}

	for i, expected := range expRespSlice {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(expected)
			}))
			defer svr.Close()
			client := &http.Client{}
			resp, err := getWeather(svr.URL, client)
			if err != nil {
				t.Errorf("expected err to be nil got %v", err)
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if string(body) != string(expected) {
				t.Errorf("expected response body to be %s got %s", expected, body)
			}
		})
	}
}
