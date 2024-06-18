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

func Scheduler(p *kafka.Producer, jobs []models.CronJob, mapTopics map[string]string, duration time.Duration) error {

	c := cron.New(cron.WithSeconds())

	for _, job := range jobs {
		topic := mapTopics[job.Cluster]
		if topic == "" {
			log.Printf("Warning: Jobs has an empty topic for cluster %s", job.Cluster)
			continue
		}

		err := scheduleJob(c, p, job, topic)
		if err != nil {
			return fmt.Errorf("error scheduling job: %w", err)
		}
	}

	c.Start()

	time.Sleep(duration)

	c.Stop()
	return nil
}

func scheduleJob(c *cron.Cron, p *kafka.Producer, job models.CronJob, topic string) error {
	_, err := c.AddFunc(job.Schedule, func() {
		err := producer.ProduceMessage(p, topic, job)
		if err != nil {
			fmt.Printf("error producing message: %v", err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}
