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
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheus.MustRegister(models.CronJobLatency, models.CronJobCount, models.CronJobErrorCount, models.CronJobRetryCount)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":2112", nil)
	}()

	ts, brokers, cluster, retry := utils.Args()
	topics := strings.Split(ts, ",")

	// Map the topics to the clusters
	mapTopics := map[string]string{
		"cluster-a": topics[0],
		"cluster-b": topics[1],
	}

	topic := mapTopics[cluster]

	if retry {
		topic = topic + "-retry"
	}

	consumer, err := initialiseConsumer(brokers, topic)
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

func initialiseConsumer(brokers string, topic string) (*kafka.Consumer, error) {
	cons, err := consumer.Consumer(brokers)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	err = cons.Subscribe(topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)

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
				fmt.Printf("Job %s failed: %v\n", job.Id, err)
				if job.Retries > 0 {
					job.Retries--
					RetryJob(p, job, job.RetryTopic)
					time.Sleep(time.Duration(job.RetryInterval) * time.Second)
					models.CronJobRetryCount.WithLabelValues(job.Cluster).Inc()
				} else {
					fmt.Printf("No more retries for job: %v\n", job.Id)
					models.CronJobErrorCount.WithLabelValues(job.Cluster).Inc()
					break
				}
			} else {
				fmt.Printf("Job %s executed successfully\n", job.Id)
				continue
			}

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func RetryJob(p *kafka.Producer, job models.CronJob, topic string) {
	fmt.Printf("Retrying job %s\n", job.Id)
	err := producer.ProduceMessage(p, topic, job)
	if err != nil {
		fmt.Printf("error producing message: %v", err)
	}
}
