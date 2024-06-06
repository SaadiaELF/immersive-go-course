package main

import (
	"encoding/json"
	"fmt"
	"kafka-cron/consumer/internal/consumer"
	"kafka-cron/consumer/internal/executor"
	"kafka-cron/pkg/models"
	"kafka-cron/utils"
	"os"
)

func main() {
	topic, brokers := utils.Args()
	consumer, err := consumer.Consumer(brokers)
	if err != nil {
		fmt.Printf("failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		fmt.Printf("failed to subscribe to topic: %s\n", err)
		os.Exit(1)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: ", msg.TopicPartition)
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
