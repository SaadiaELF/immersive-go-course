package main

import (
	"encoding/json"
	"io"
	"kafka-cron/cmd/producer/pkg/models"
	"log"
	"os"
)

// Read and parse the configuration file
func readConfig() ([]byte, error) {

	path := "../../cron-config.json"

	// Open the configuration file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func parseConfig(data []byte) ([]models.CronJob, error) {
	var jobs []models.CronJob
	err := json.Unmarshal(data, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil

}

func main() {
	data, err := readConfig()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}

	_, err = parseConfig(data)
	if err != nil {
		log.Fatalf("error parsing the configuration file: %v", err)
	}

}
