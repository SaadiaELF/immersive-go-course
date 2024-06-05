package main

import (
	"fmt"
	"kafka-cron/producer/internal/producer"
	"kafka-cron/producer/internal/scheduler"
	"kafka-cron/utils"
	"log"
)

func main() {
	// Read and parse the configuration file
	data, err := utils.ReadConfig()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}

	jobs, err := utils.ParseConfig(data)
	if err != nil {
		log.Fatalf("error parsing the configuration file: %v", err)
	}

	topic, brokers := utils.Args()

	// Create the Kafka producer
	p, err := producer.Producer(brokers)
	if err != nil {
		fmt.Printf("Error creating Kafka producer: %v", err)
	}

	// Create the Kafka topic
	err = producer.CreateTopic(p, topic)
	if err != nil {
		fmt.Printf("Error creating Kafka topic: %v", err)
	}
	scheduler.Scheduler(p, jobs, topic)

}
