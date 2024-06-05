package main

import (
	"encoding/json"
	"fmt"
	"io"
	"kafka-cron/cmd/producer/pkg/models"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
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

func scheduleJobs(jobs []models.CronJob) {
	c := cron.New(cron.WithSeconds())

	// Schedule the jobs
	for _, job := range jobs {
		_, err := c.AddFunc(job.Schedule, func() {
			err := executeJob(job)
			if err != nil {
				fmt.Println("error executing job:", job.Name)
			}
		})
		if err != nil {
			fmt.Println("error scheduling job:", job.Name)
		}
	}

	// Start the Cron job scheduler
	c.Start()

	// Wait for the Cron job to run
	time.Sleep(5 * time.Minute)

	// Stop the Cron job scheduler
	c.Stop()

}

func executeJob(job models.CronJob) error {
	cmd := exec.Command("sh", "-c", job.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Executing job:", job.Name)
	return cmd.Run()
}

func main() {
	data, err := readConfig()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}

	jobs, err := parseConfig(data)
	if err != nil {
		log.Fatalf("error parsing the configuration file: %v", err)
	}

	scheduleJobs(jobs)

}
