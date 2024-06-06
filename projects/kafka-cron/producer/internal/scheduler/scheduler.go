package scheduler

import (
	"fmt"
	"kafka-cron/pkg/models"
	"kafka-cron/producer/internal/producer"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/robfig/cron/v3"
)

func Scheduler(p *kafka.Producer, jobs []models.CronJob, topic string, duration time.Duration) error {
	c := cron.New(cron.WithSeconds())

	// Schedule the jobs
	for _, job := range jobs {
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
