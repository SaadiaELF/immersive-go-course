package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
)

func TestPostImage(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	type test struct {
		name               string
		image              string
		imageUrl           string
		expectedStatusCode int
		expectSelect       bool
		expectInsert       bool
	}

	testCases := []test{
		{
			name:               "empty image",
			image:              "",
			imageUrl:           "/image",
			expectedStatusCode: http.StatusBadRequest,
			expectSelect:       false,
			expectInsert:       false,
		},
		{
			name:               "existing image url",
			image:              `{"title": "image", "url": "http://example.com/image", "alt_text": "img"}`,
			imageUrl:           "http://example.com/image",
			expectedStatusCode: http.StatusConflict,
			expectSelect:       true,
			expectInsert:       false,
		},
		{
			name:               "insert image",
			image:              `{"title": "image", "url": "http://example.com/image", "alt_text": "img"}`,
			imageUrl:           "http://example.com/image",
			expectedStatusCode: http.StatusCreated,
			expectSelect:       false,
			expectInsert:       true,
		},
	}

	for _, tc := range testCases {

		if tc.expectSelect {
			mock.NewRows([]string{"url"}).AddRow("http://example.com/image")
			mock.ExpectQuery(`SELECT url FROM public.images WHERE url = \$1`).
				WithArgs(tc.imageUrl).
				WillReturnRows(mock.NewRows([]string{"url"}).AddRow("http://example.com/image"))
		}
		if tc.expectInsert {
			query := `INSERT INTO public.images`
			mock.ExpectExec(query).
				WithArgs("image", "http://example.com/image", "img").
				WillReturnResult(pgxmock.NewResult("INSERT", 1))
		}

		req, _ := http.NewRequest(http.MethodPost, tc.imageUrl, strings.NewReader(tc.image))
		resp := httptest.NewRecorder()
		postImage(mock, resp, req)
		if tc.expectedStatusCode != resp.Code {
			t.Errorf("testcase %s: got %v, want %v", tc.name, resp.Code, tc.expectedStatusCode)
		}
	}
}

func TestFetchImages(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	query := `SELECT title, url, alt_text FROM public.images LIMIT \$1`
	mock.ExpectQuery(query).
		WithArgs(10).
		WillReturnRows(mock.NewRows([]string{"title", "url", "alt_text"}).
			AddRow("image", "http://example.com/image", "img"))

	_, err = fetchImages(mock, 10)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

}
