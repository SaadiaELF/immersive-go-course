package main

import (
	"encoding/json"
	"fmt"
	"kafka-cron/consumer/internal/consumer"
	"kafka-cron/consumer/internal/executor"
	"kafka-cron/consumer/internal/producer"
	"kafka-cron/consumer/utils"
	"kafka-cron/pkg/models"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheus.MustRegister(models.CronJobLatency)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":2112", nil)
	}()

	topic1, topic2, brokers, cluster, retry := utils.Args()

	var topic string
	switch cluster {
	case "cluster-a":
		topic = topic1
		if retry {
			topic = topic + "-retry"
		}
	case "cluster-b":
		topic = topic2
		if retry {
			topic = topic + "-retry"
		}
	default:
		fmt.Printf("Invalid cluster specified: %s. Use 'cluster-a' or 'cluster-b'.\n", cluster)
		return
	}

	consumer, err := initializeConsumer(brokers, topic)
	if err != nil {
		fmt.Printf("failed to initialise consumer: %s\n", err)
	}

	// Create the Kafka producer
	p, err := producer.Producer(brokers)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)

	}
	done := make(chan bool)

	go func() {
		processMessages(p, consumer)
	}()

	done <- true
}

func initializeConsumer(brokers string, topic string) (*kafka.Consumer, error) {
	cons, err := consumer.Consumer(brokers)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %s", err)
	}
	err = cons.Subscribe(topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %s", err)

	}

	return cons, nil
}

func processMessages(p *kafka.Producer, con *kafka.Consumer) {
	for {
		msg, err := con.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message for %s on %s: \n", *msg.TopicPartition.Topic, msg.TopicPartition)
			var job models.CronJob
			err = json.Unmarshal(msg.Value, &job)
			if err != nil {
				fmt.Printf("failed to unmarshal message: %v\n", err)
				continue
			}
			err = executor.Execute(job)

			if err != nil {
				fmt.Printf("Job failed: %v\n", err)
				if job.Retries > 0 {
					job.Retries--
					RetryJob(p, job, job.RetryTopic)
					time.Sleep(time.Duration(job.RetryInterval) * time.Second)
				}
				if job.Retries == 0 {
					fmt.Printf("No more retries for job: %v\n", job.Id)
				}
			}

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func RetryJob(p *kafka.Producer, job models.CronJob, topic string) {
	fmt.Println("Retrying job")
	err := producer.ProduceMessage(p, topic, job)
	if err != nil {
		fmt.Printf("error producing message: %v", err)
	}
}
