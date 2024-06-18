package main

import (
	"encoding/json"
	"fmt"
	"kafka-cron/consumer/internal/consumer"
	"kafka-cron/consumer/internal/executor"
	"kafka-cron/consumer/utils"
	"kafka-cron/pkg/models"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	ts, brokers, cluster := utils.Args()
	topics := strings.Split(ts, ",")

	// Map the topics to the clusters
	mapTopics := map[string]string{
		"cluster-a": topics[0],
		"cluster-b": topics[1],
	}

	topic := mapTopics[cluster]

	consumer, err := initialiseConsumer(brokers, topic)
	if err != nil {
		fmt.Printf("failed to initialise consumer: %s\n", err)
	}

	done := make(chan bool)

	go func() {
		processMessages(consumer)
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

func processMessages(con *kafka.Consumer) {
	for {
		msg, err := con.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message for %s on %s: ", *msg.TopicPartition.Topic, msg.TopicPartition)
			var job models.CronJob
			err = json.Unmarshal(msg.Value, &job)
			if err != nil {
				fmt.Printf("failed to unmarshal message: %v\n", err)
				continue
			}
			err = executor.Execute(job)
			if err != nil {
				fmt.Printf("failed to execute job: %v\n", err)
			}
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
