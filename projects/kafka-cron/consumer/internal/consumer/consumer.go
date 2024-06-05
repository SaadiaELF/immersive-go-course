package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Consumer() (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
