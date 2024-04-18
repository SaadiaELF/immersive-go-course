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
		selectResultRow    *pgxmock.Rows
		expectInsert       bool
		// expectInsertValues []
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
			image:              `{"title": "image", "url": "/image", "alt_text": "img"}`,
			imageUrl:           "/image",
			expectedStatusCode: http.StatusConflict,
			expectSelect:       true,
			selectResultRow:    mock.NewRows([]string{"url"}).AddRow("/image"),
			expectInsert:       false,
		},
	}

	for _, tc := range testCases {

		if tc.expectSelect {
			query := `SELECT url FROM public.images WHERE url = $1`
			mock.ExpectQuery(query).
				WithArgs(tc.imageUrl).
				WillReturnRows(tc.selectResultRow)
		}
		// if tc.expectInsert {
		// TODO: Implement the expectInsert branch.
		// }

		req, _ := http.NewRequest(http.MethodPost, tc.imageUrl, strings.NewReader(tc.image))
		resp := httptest.NewRecorder()
		postImage(resp, req)
		if tc.expectedStatusCode != resp.Code {
			t.Errorf("testcase %s: got %v, want %v", tc.name, tc.expectedStatusCode, resp.Code)
		}
	}
}
