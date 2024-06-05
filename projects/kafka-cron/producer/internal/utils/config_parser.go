package utils

import (
	"encoding/json"
	"io"
	"kafka-cron/producer/pkg/models"
	"os"
)

func ReadConfig() ([]byte, error) {

	path := "./config/cron-config.json"

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
func ParseConfig(data []byte) ([]models.CronJob, error) {
	var jobs []models.CronJob
	err := json.Unmarshal(data, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil

}
