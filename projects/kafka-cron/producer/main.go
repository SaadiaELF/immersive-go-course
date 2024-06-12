package main

import (
	"kafka-cron/producer/config"
	"kafka-cron/producer/internal/producer"
	"kafka-cron/producer/internal/scheduler"
	"kafka-cron/producer/utils"
	"log"
	"time"
)

func main() {
	topic1, topic2, brokers := utils.Args()

	// Read and parse the configuration file
	data, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("error reading the configuration file: %v", err)
	}

	jobs, err := config.ParseConfig(data)
	if err != nil {
		log.Fatalf("error parsing the configuration file: %v", err)
	}

	// Create the Kafka producer
	p, err := producer.Producer(brokers)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)

	}

	// Create the Kafka topics
	if topic1 != "" {
		err = producer.CreateTopic(p, topic1)
		if err != nil {
			log.Printf("Error creating Kafka topic1: %v", err)
		}
		//create retry topic
		err = producer.CreateTopic(p, topic1+"-retry")
		if err != nil {
			log.Printf("Error creating Kafka topic1-retry: %v", err)
		}
	}

	if topic2 != "" {
		err = producer.CreateTopic(p, topic2)
		if err != nil {
			log.Printf("Error creating Kafka topic2: %v", err)
		}
		//create retry topic
		err = producer.CreateTopic(p, topic2+"-retry")
		if err != nil {
			log.Printf("Error creating Kafka topic2-retry: %v", err)
		}
	}

	done := make(chan bool)
	go func() {

		err := scheduler.Scheduler(p, jobs, topic1, topic2, 5*time.Minute)
		if err != nil {
			log.Printf("Error scheduling jobs for cluster: %v", err)
		}
	}()
	done <- true

	// Close the producer
	p.Close()
}
