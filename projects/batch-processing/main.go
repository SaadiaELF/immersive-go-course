package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
		log.Printf("error: Could not read csv file: %v\n", err)
	}

	outputRecords := [][]string{{"url", "input", "output", "s3url"}}

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

		// Download the image
		inputFilename := fmt.Sprintf("%s/img-0%v.jpg", *inputFilepath, i)
		err := DownloadImage(inputFilename, record[0])
		if err != nil {
			log.Printf("error downloading: %v\n", err)
		}

		// Convert the image to grayscale
		outputFilename := fmt.Sprintf("%s/img-0%v.jpg", *outputFilepath, i)
		err = c.Grayscale(inputFilename, outputFilename)
		if err != nil {
			log.Printf("error converting image: %v\n", err)
		}

		//Upload the images to the aws s3 bucket
		s3url, err := UploadImage(outputFilename)
		if err != nil {
			log.Printf("error uploading image: %v\n", err)
		}
		outputRecords = append(outputRecords, []string{record[0], inputFilename, outputFilename, s3url})
	}
	// Create a CSV file with the output records
	log.Println("Creating output CSV file ... ")
	_, err = CreateCSVFile(outputRecords)
	if err != nil {
		log.Printf("error creating output csv file: %v\n", err)
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

	if len(records) == 0 {
		return nil, fmt.Errorf("no records found in the csv file")
	}

	if len(records[0]) > 1 {
		return records, fmt.Errorf("more than one column is found in the csv file")
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get the image: %v", resp.Status)
	}

	// Check if it's an image
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("invalid image type: %v", contentType)
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

func UploadImage(filename string) (string, error) {
	// Get the AWS region and role ARN from the environment
	region := os.Getenv("AWS_REGION")
	awsRoleArn := os.Getenv("AWS_ROLE_ARN")
	bucket := os.Getenv("AWS_BUCKET")
	if region == "" || awsRoleArn == "" {
		return "", fmt.Errorf("AWS_REGION and AWS_ROLE_ARN environment variables must be set")
	}
	// Set up S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", fmt.Errorf("error creating session: %v\n", err)
	}

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the ARN.
	creds := stscreds.NewCredentials(sess, awsRoleArn)

	// Create service client value configured for credentials
	// from assumed role.
	svc := s3.New(sess, &aws.Config{Credentials: creds, Endpoint: aws.String("s3." + region + ".amazonaws.com")})

	// // Upload the image to the s3 bucket
	bufBytes := getFileBytes(filename)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(bufBytes),
	})
	if err != nil {
		return "", err
	}
	// Construct the URL of the uploaded image
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, filename)
	return url, nil
}

func getFileBytes(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("could not open file: %v", err)
	}
	defer file.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		log.Printf("could not copy file: %v", err)
	}
	return buf.Bytes()
}

func CreateCSVFile(records [][]string) (string, error) {
	// Create a temporary file
	file, err := os.Create("./outputs/output.csv")
	if err != nil {
		return "", fmt.Errorf("could not create csv file: %v", err)
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
