package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/gographics/imagick.v2/imagick"
)

type ConvertImageCommand func(args []string) (*imagick.ImageCommandResult, error)

type Converter struct {
	cmd ConvertImageCommand
}

func main() {
	// Accept --input and --output arguments for the images
	inputFilepath := flag.String("input", "", "A path to an image to be processed")
	outputFilepath := flag.String("output", "", "A path to where the processed image should be written")
	flag.Parse()

	// Ensure that both flags were set
	if *inputFilepath == "" || *outputFilepath == "" {
		flag.Usage()
		os.Exit(1)
	}
	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}
	// Read the CSV file
	log.Println("Reading input CSV file ... ")
	records, err := ReadCSV("./inputs/inputs.csv")
	if err != nil {
		log.Fatalf("error: Could not read csv file: %v\n", err)
	}
	if len(records) == 0 {
		log.Fatalln("no records found in the csv file")
	}
	if len(records[0]) > 1 {
		log.Println("more than one column is found in the csv file")
	}
	// Download the images and process them
	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Log what we're going to do
	log.Printf("processing: %q to %q\n", *inputFilepath, *outputFilepath)

	for i, record := range records {
		// Check if the first row is the header
		if i == 0 {
			if record[0] == "url" {
				continue
			} else {
				log.Fatalln("no url header found in the csv file")
			}
		}

		// Check if the url is valid
		if !IsValidURL(record[0]) {
			log.Printf("invalid url found in the csv file: %v\n", record[0])
		}

		// Download the image
		filename := fmt.Sprintf("%s/img-0%v.jpg", *inputFilepath, i)
		err := DownloadImage(filename, record[0])
		if err != nil {
			log.Printf("error downloading: %v\n", err)
		}

		// Convert the image to grayscale
		dest := fmt.Sprintf("%s/img-0%v.jpg", *outputFilepath, i)
		err = c.Grayscale(filename, dest)
		if err != nil {
			log.Printf("error converting image: %v\n", err)
		}
	}

	// Log what we did
	log.Printf("processed: %q to %q\n", *inputFilepath, *outputFilepath)
}

func ReadCSV(filename string) (records [][]string, err error) {
	// Open the file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the CSV file
	r := csv.NewReader(f)
	records, err = r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
func IsValidURL(imgUrl string) bool {
	_, err := url.ParseRequestURI(imgUrl)
	if err != nil {
		return false
	}
	return true
}

func DownloadImage(filepath string, url string) error {
	// Create empty file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the image
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the image to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (c *Converter) Grayscale(inputFilepath string, outputFilepath string) error {
	// Convert the image to grayscale using imagemagick
	// We are directly calling the convert command
	_, err := c.cmd([]string{
		"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
	})
	return err
}
