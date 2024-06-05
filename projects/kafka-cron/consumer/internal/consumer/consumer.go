package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Consumer(brokers string) (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "myGroup",
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}
