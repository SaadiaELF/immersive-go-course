package producer

import (
	"context"
	"fmt"
	"kafka-cron/pkg/models"
	"time"

	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
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

func CreateTopic(p *kafka.Producer, topic string) error {
	a, err := kafka.NewAdminClientFromProducer(p)
	if err != nil {
		return fmt.Errorf("failed to create new admin client from producer: %w", err)
	}
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create topics on cluster.
	// Set Admin options to wait up to 60s for the operation to finish on the remote cluster
	maxDur := time.Second * 60
	results, err := a.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{
			{
				Topic:             topic,
				NumPartitions:     2,
				ReplicationFactor: 1,
			},
		},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		return fmt.Errorf("admin Client request error: %w", err)
	}
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			return fmt.Errorf("failed to create topic: %w", result.Error)
		}
		fmt.Printf("%v\n", result)
	}
	a.Close()
	return nil
}

func ProduceMessage(p *kafka.Producer, topic string, job models.CronJob) error {
	deliveryChan := make(chan kafka.Event)
	message, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(uuid.New().String()),
		Value:          []byte(message),
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}
	fmt.Printf("delivered %s to topic %s [%d] at offset %v\n", m.Key,
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	close(deliveryChan)
	return nil
}
