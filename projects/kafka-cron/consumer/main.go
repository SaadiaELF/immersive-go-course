package main

import (
	"encoding/json"
	"fmt"
	"kafka-cron/consumer/internal/consumer"
	"kafka-cron/consumer/internal/executor"
	"kafka-cron/pkg/models"
	"kafka-cron/utils"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	topic1, topic2, brokers := utils.Args()
	consumer1, err := initializeConsumer(brokers, topic1)
	if err != nil {
		fmt.Printf("failed to initialise consumer: %s\n", err)
	}
	consumer2, err := initializeConsumer(brokers, topic2)
	if err != nil {
		fmt.Printf("failed to initialise consumer: %s\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		processMessages(consumer1)
	}()

	go func() {
		defer wg.Done()
		processMessages(consumer2)
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
				fmt.Printf("failed to unmarshal message: %s\n", err)
				continue
			}
			executor.Executor(job)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
