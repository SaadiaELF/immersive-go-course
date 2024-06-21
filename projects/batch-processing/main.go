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

type ImageUploader struct {
	region     string
	awsRoleArn string
	bucket     string
	s3         *s3.S3
}

func main() {
	inputFilepath := flag.String("input", "", "A path to a CSV file with image URLs to process")
	outputFilepath := flag.String("output", "", "A path to where the processed records should be written")
	failedFilepath := flag.String("output-failed", "", "A path to where the failed records should be written")
	flag.Parse()

	if *inputFilepath == "" || *outputFilepath == "" || *failedFilepath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Read and validate environment variables
	region := os.Getenv("AWS_REGION")
	roleArn := os.Getenv("AWS_ROLE_ARN")
	bucket := os.Getenv("AWS_BUCKET")
	if region == "" || roleArn == "" || bucket == "" {
		log.Println("AWS_REGION, AWS_ROLE_ARN and AWS_BUCKET environment variables must be set")
	}

	// Initialise the ImageUploader
	uploader, err := NewImageUploader(region, roleArn, bucket)
	if err != nil {
		log.Printf("could not create image uploader: %v", err)
	}

	// Build a Converter struct that will use imagick
	c := &Converter{
		cmd: imagick.ConvertImageCommand,
	}

	// Log what we're going to do
	log.Printf("processing: %q to %q\n", *inputFilepath, *outputFilepath)

	log.Println("reading input CSV file ... ")
	inputsFilepath := fmt.Sprintf("./inputs/%s", *inputFilepath)
	records, err := ReadCSV(inputsFilepath)
	if err != nil {
		log.Printf("error: Could not read csv file: %v\n", err)
	}

	// Create csv file for output records
	log.Println("creating output CSV file ... ")
	outputsFilepath := fmt.Sprintf("./outputs/%s", *outputFilepath)
	outputFile, err := os.Create(outputsFilepath)
	if err != nil {
		log.Fatalf("could not create csv file: %v", err)
	}
	defer outputFile.Close()
	outputWriter := csv.NewWriter(outputFile)
	defer outputWriter.Flush()

	// Create csv file for failed records
	log.Println("Creating failed CSV file ... ")
	failedPath := fmt.Sprintf("./outputs/%s", *failedFilepath)
	failedFile, err := os.Create(failedPath)
	if err != nil {
		log.Fatalf("could not create csv file: %v", err)
	}
	defer failedFile.Close()
	failedWriter := csv.NewWriter(failedFile)
	defer failedWriter.Flush()

	WriteHeader(outputWriter, []string{"url", "input", "output", "s3url"})
	WriteHeader(failedWriter, []string{"url"})

	// Download the images and process them
	// Set up imagemagick
	imagick.Initialize()
	defer imagick.Terminate()

	ProcessImage(records, c, uploader, outputWriter, failedWriter)

	log.Printf("processed: %q to %q\n", *inputFilepath, *outputFilepath)
}

func ReadCSV(filepath string) (records [][]string, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err = r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no records found in the csv file")
	}
	if records[0][0] != "url" {
		return nil, fmt.Errorf("no url header found in the csv file")
	}
	if len(records[0]) > 1 {
		return records, fmt.Errorf("more than one column is found in the csv file")
	}

	return records, nil
}

func WriteHeader(writer *csv.Writer, header []string) error {
	err := writer.Write(header)
	if err != nil {
		return fmt.Errorf("could not write header to csv file: %v", err)
	}
	return nil
}

func ProcessImage(records [][]string, c *Converter, uploader *ImageUploader, outputWriter, failedWriter *csv.Writer) {
	for i := 1; i < len(records); i++ {
		// Download the image
		inputFilename, err := DownloadImage(records[i][0])
		if err != nil {
			log.Printf("error downloading: %v\n", err)
			writeFailedRecord(failedWriter, records[i][0])
			continue
		}

		// Convert the image to grayscale
		outputFilename := fmt.Sprintf("/tmp/img-0%v.jpg", i)
		err = c.Grayscale(inputFilename, outputFilename)
		if err != nil {
			log.Printf("error converting image: %v\n", err)
			writeFailedRecord(failedWriter, records[i][0])
			continue
		}

		//Upload the images to the aws s3 bucket if no error
		s3url, err := uploader.UploadImage(outputFilename)
		if err != nil {
			log.Printf("error uploading image: %v\n", err)
			writeFailedRecord(failedWriter, records[i][0])
		}

		// Write the record to the output csv file
		outputRecord := []string{records[i][0], inputFilename, outputFilename, s3url}
		if err := outputWriter.Write(outputRecord); err != nil {
			log.Printf("could not write record to csv file: %v", err)
		}
	}
}
func writeFailedRecord(writer *csv.Writer, url string) {
	if err := writer.Write([]string{url}); err != nil {
		log.Printf("could not write record to failed csv file: %v", err)
	}
}

func DownloadImage(url string) (string, error) {
	// Create empty file
	out, err := os.CreateTemp("", "img-*.jpg")
	if err != nil {
		return "", err
	}
	defer out.Close()
	filename := out.Name()

	// Get the image
	resp, err := http.Get(url)
	if err != nil {
		return filename, err
	}
	if resp.StatusCode != http.StatusOK {
		return filename, fmt.Errorf("failed to get the image: %v", resp.Status)
	}

	// Check if it's an image
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return filename, fmt.Errorf("invalid image type: %v", contentType)
	}

	defer resp.Body.Close()

	// Write the image to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return filename, err
	}
	return filename, nil
}

func (c *Converter) Grayscale(inputFilepath string, outputFilepath string) error {
	// Convert the image to grayscale using imagemagick
	// We are directly calling the convert command
	_, err := c.cmd([]string{
		"convert", inputFilepath, "-set", "colorspace", "Gray", outputFilepath,
	})
	return err
}

func NewImageUploader(region, awsRoleArn, bucket string) (*ImageUploader, error) {
	// Set up S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating session: %v\n", err)
	}

	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the ARN.
	creds := stscreds.NewCredentials(sess, awsRoleArn)

	// Create service client value configured for credentials
	// from assumed role.
	svc := s3.New(sess, &aws.Config{Credentials: creds, Endpoint: aws.String("s3." + region + ".amazonaws.com")})

	return &ImageUploader{
		region:     region,
		awsRoleArn: awsRoleArn,
		bucket:     bucket,
		s3:         svc,
	}, nil

}

func (u *ImageUploader) UploadImage(filename string) (string, error) {
	// Upload the image to the S3 bucket
	bufBytes := getFileBytes(filename)
	_, err := u.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(bufBytes),
	})
	if err != nil {
		return "", err
	}
	// Construct the URL of the uploaded image
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, filename)
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
