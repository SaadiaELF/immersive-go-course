package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func TestGrayscaleMockError(t *testing.T) {
	c := &Converter{
		cmd: func(args []string) (*imagick.ImageCommandResult, error) {
			return nil, errors.New("not implemented")
		},
	}

	err := c.Grayscale("input.jpg", "output.jpg")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGrayscaleMockCall(t *testing.T) {
	var args []string
	expected := []string{"convert", "input.jpg", "-set", "colorspace", "Gray", "output.jpg"}
	c := &Converter{
		cmd: func(a []string) (*imagick.ImageCommandResult, error) {
			args = a
			return &imagick.ImageCommandResult{
				Info: nil,
				Meta: "",
			}, nil
		},
	}

	err := c.Grayscale("input.jpg", "output.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, args) {
		t.Fatalf("incorrect arguments: expected %v, got %v", expected, args)
	}
}

func CreateTempCSVFile(records [][]string) (string, error) {
	// Create a temporary file
	file, err := os.CreateTemp("", "input-*.csv")
	if err != nil {
		return "", fmt.Errorf("could not create temp file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer and write the records to the file
	writer := csv.NewWriter(file)
	err = writer.WriteAll(records)
	if err != nil {
		return "", fmt.Errorf("could not write records to CSV: %v", err)
	}
	return file.Name(), nil
}

func TestReadCSV(t *testing.T) {
	testcases := []struct {
		name    string
		records [][]string
	}{
		{
			name:    "empty file",
			records: [][]string{},
		},
		{
			name: "two columns",
			records: [][]string{
				{"url", "name"},
				{"foo/bar", "bar"},
			},
		},
		{
			name: "valid file",
			records: [][]string{
				{"url"},
				{"foo/bar"},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			filepath, err := CreateTempCSVFile(tc.records)
			if err != nil {
				t.Fatal(err)
			}
			records, err := ReadCSV(filepath)
			if tc.name != "valid file" && err == nil {
				t.Fatalf("expected error, got %v", records)
			}
			if tc.name == "valid file" && err != nil {
				t.Fatalf("expected records, got error: %v", err)
			}

			defer os.Remove(filepath)
		})
	}

}

func CreateTempImageFile() (string, error) {
	// Create a temporary file
	file, err := os.CreateTemp("", "image-*.jpg")
	if err != nil {
		return "", fmt.Errorf("could not create temp file: %v", err)
	}
	return file.Name(), nil
}
func TestDownloadImage(t *testing.T) {
	testcases := []struct {
		name string
		url  string
	}{
		{
			name: "invalid url non-200",
			url:  "https://example.com/image.jpg",
		},
		{
			name: "invalid url non-image",
			url:  "https://google.com",
		},

		{
			name: "invalid url",
			url:  "/invalid/url",
		},
		{
			name: "valid url",
			url:  "https://upload.wikimedia.org/wikipedia/commons/1/11/Test-Logo.svg",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			filepath, err := CreateTempImageFile()
			if err != nil {
				t.Fatal(err)
			}
			err = DownloadImage(filepath, tc.url)
			if tc.name != "valid url" && err == nil {
				t.Fatalf("expected error, got %v", filepath)
			}
			if tc.name == "valid url" && err != nil {
				t.Fatalf("expected %s, got error: %v", filepath, err)
			}

			defer os.Remove(filepath)
		})
	}
}
