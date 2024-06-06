package main

import (
	"kafka-cron/pkg/models"
	"kafka-cron/producer/internal/producer"
	"kafka-cron/producer/internal/scheduler"
	"kafka-cron/producer/utils"
	"log"
	"sync"
	"time"
)

func main() {
	topic1, topic2, brokers, cluster := utils.Args()

	jobs, err := utils.ParseConfig(cluster)
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

	// Schedule the jobs
	clusterAJobs := []models.CronJob{}
	clusterBJobs := []models.CronJob{}

	for _, job := range jobs {
		switch job.Cluster {
		case "cluster-a":
			clusterAJobs = append(clusterAJobs, job)
		case "cluster-b":
			clusterBJobs = append(clusterBJobs, job)
		default:
			log.Printf("Warning: Jobs has an unknown cluster %s", job.Cluster)
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := scheduler.Scheduler(p, clusterAJobs, topic1, 5*time.Minute)
		if err != nil {
			log.Printf("Error scheduling jobs for cluster a: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err = scheduler.Scheduler(p, clusterBJobs, topic2, 5*time.Minute)
		if err != nil {
			log.Printf("Error scheduling jobs for cluster b:  %v", err)
		}
	}()

	wg.Wait()

	// Close the producer
	p.Close()

}
