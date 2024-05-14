package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/gographics/imagick.v2/imagick"
)

type ConvertImageCommand func(args []string) (*imagick.ImageCommandResult, error)

type Converter struct {
	cmd ConvertImageCommand
}

func (c *Converter) Grayscale(inputFilepath string, outputFilepath string) error {
	// Convert the image to grayscale using imagemagick
	// We are directly calling the convert command
	_, err := c.cmd([]string{
		"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
	})
	return err
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

	// Read the CSV file
	log.Println("Reading input CSV file ... ")
	records, _ := ReadCSV("./inputs/inputs.csv")

	// Download the images and process them
	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Log what we're going to do
	log.Printf("processing: %q to %q\n", *inputFilepath, *outputFilepath)

	for i, record := range records {
		// Skip the header
		if i == 0 {
			continue
		}

		filename := fmt.Sprintf("%s/img-0%v.jpg", *inputFilepath, i)
		fmt.Println(filename)
		err := DownloadImage(filename, record[0])
		if err != nil {
			log.Printf("error: %v\n", err)
		}

		// Build a Converter struct that will use imagick
		c := &Converter{
			cmd: imagick.ConvertImageCommand,
		}

		// Do the conversion!
		dest := fmt.Sprintf("%s/img-0%v.jpg", *outputFilepath, i)
		err = c.Grayscale(filename, dest)
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}

	// Log what we did
	log.Printf("processed: %q to %q\n", *inputFilepath, *outputFilepath)
}
