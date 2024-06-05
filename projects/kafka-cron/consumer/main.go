package main

import (
	"fmt"
	"kafka-cron/consumer/internal/consumer"
	"os"
)

func main() {

	consumer, err := consumer.Consumer()
	if err != nil {
		fmt.Printf("failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	err = consumer.Subscribe("cron-topic", nil)
	if err != nil {
		fmt.Printf("failed to subscribe to topic: %s\n", err)
		os.Exit(1)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}
