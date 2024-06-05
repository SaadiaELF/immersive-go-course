package main

import (
	"kafka-cron/producer/internal/scheduler"
	"kafka-cron/producer/internal/utils"
	"log"
)

func main() {
	data, err := utils.ReadConfig()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}

	jobs, err := utils.ParseConfig(data)
	if err != nil {
		log.Fatalf("error parsing the configuration file: %v", err)
	}

	scheduler.Scheduler(jobs)

}
