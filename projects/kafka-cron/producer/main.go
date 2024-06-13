package main

import (
	"kafka-cron/producer/config"
	"kafka-cron/producer/internal/producer"
	"kafka-cron/producer/internal/scheduler"
	"kafka-cron/producer/utils"
	"log"
	"strings"
	"time"
)

func main() {

	ts, cs, brokers := utils.Args()
	topics := strings.Split(ts, ",")
	clusters := strings.Split(cs, ",")

	// Map the topics to the clusters
	if len(clusters) != len(topics) {
		log.Fatalf("Number of clusters and topics do not match")
	}
	mapTopics := make(map[string]string)
	for i, cluster := range clusters {
		mapTopics[cluster] = topics[i]
	}

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
	for _, topic := range topics {
		err = producer.CreateTopic(p, topic)
		if err != nil {
			log.Printf("Error creating Kafka topic: %v", err)
		}
		//create retry topic
		err = producer.CreateTopic(p, topic+"-retry")
		if err != nil {
			log.Printf("Error creating Kafka topic2-retry: %v", err)
		}
	}

	done := make(chan bool)
	go func() {
		err := scheduler.Scheduler(p, jobs, mapTopics, 5*time.Minute)
		if err != nil {
			log.Printf("Error scheduling jobs for cluster: %v", err)
		}
	}()
	done <- true

	// Close the producer
	p.Close()
}
