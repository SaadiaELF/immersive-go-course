package scheduler

import (
	"fmt"
	"kafka-cron/pkg/models"
	"kafka-cron/producer/internal/producer"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/robfig/cron/v3"
)

func Scheduler(p *kafka.Producer, jobs []models.CronJob, topic1 string, topic2 string, duration time.Duration) error {

	c := cron.New(cron.WithSeconds())

	// Schedule the jobs
	for _, job := range jobs {
		var topic string
		switch job.Cluster {
		case "cluster-a":
			topic = topic1
		case "cluster-b":
			topic = topic2
		default:
			log.Printf("Warning: Jobs has an unknown cluster %s", job.Cluster)
		}
		_, err := c.AddFunc(job.Schedule, func() {
			err := producer.ProduceMessage(p, topic, job)
			if err != nil {
				fmt.Printf("error producing message: %v", err)
			}
		})
		if err != nil {
			return fmt.Errorf("error scheduling job: %w", err)
		}
	}

	c.Start()

	time.Sleep(duration)

	c.Stop()
	return nil
}
