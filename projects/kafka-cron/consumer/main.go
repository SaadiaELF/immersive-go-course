package main

import (
	"encoding/json"
	"fmt"
	"kafka-cron/consumer/internal/consumer"
	"kafka-cron/consumer/internal/executor"
	"kafka-cron/consumer/utils"
	"kafka-cron/pkg/models"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	topic1, topic2, brokers, cluster := utils.Args()

	var topic string
	switch cluster {
	case "cluster-a":
		topic = topic1
	case "cluster-b":
		topic = topic2
	default:
		fmt.Printf("Invalid cluster specified: %s. Use 'cluster-a' or 'cluster-b'.\n", cluster)
		return
	}

	consumer, err := initializeConsumer(brokers, topic)
	if err != nil {
		fmt.Printf("failed to initialise consumer: %s\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		processMessages(consumer)
	}()
	wg.Wait()
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
