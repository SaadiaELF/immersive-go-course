package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeather(t *testing.T) {
	expected := []byte("Today it will be sunny!")
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
		t.Errorf("expected res to be %s got %s", expected, body)
	}
}
