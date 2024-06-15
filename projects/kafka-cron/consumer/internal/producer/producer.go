package producer

import (
	"fmt"
	"kafka-cron/pkg/models"

	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Producer(brokers string) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	},
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func ProduceMessage(p *kafka.Producer, topic string, job models.CronJob) error {
	message, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value: []byte(message),
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	} else {
		fmt.Println("Produced retry job for job ", job.Id)
	}
	return nil
}
