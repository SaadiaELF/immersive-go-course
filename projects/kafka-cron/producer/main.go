package main

import (
	"kafka-cron/producer/internal/producer"
	"kafka-cron/producer/internal/scheduler"
	"kafka-cron/producer/utils"
	"log"
	"time"
)

func main() {

	topic1, topic2, brokers, _ := utils.Args()

	// Read and parse the configuration file
	data, err := utils.ReadConfig()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}

	jobs, err := utils.ParseConfig(data)
	if err != nil {
		log.Fatalf("error parsing the configuration file: %v", err)
	}

	// Create the Kafka producer
	p, err := producer.Producer(brokers)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)

	}

	// Create the Kafka topics
	err = producer.CreateTopic(p, topic1)
	if err != nil {
		log.Printf("Error creating Kafka topic1: %v", err)
	}
	err = producer.CreateTopic(p, topic2)
	if err != nil {
		log.Printf("Error creating Kafka topic2: %v", err)
	}

	done := make(chan bool)
	go func() {

		err := scheduler.Scheduler(p, jobs, topic1, topic2, 5*time.Minute)
		if err != nil {
			log.Printf("Error scheduling jobs for cluster a: %v", err)
		}
	}()
	done <- true

	// Close the producer
	p.Close()
}
