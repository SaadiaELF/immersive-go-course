package scheduler

import (
	"fmt"
	"kafka-cron/producer/internal/producer"
	"kafka-cron/producer/pkg/models"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/robfig/cron/v3"
)

func Scheduler(p *kafka.Producer, jobs []models.CronJob) error {
	c := cron.New(cron.WithSeconds())

	// Schedule the jobs
	for _, job := range jobs {
		_, err := c.AddFunc(job.Schedule, func() {
			err := producer.ProduceMessage(p, "cron-topic", job)
			if err != nil {
				fmt.Printf("error executing job: %s", job.Name)
			}
		})
		if err != nil {
			return fmt.Errorf("error executing job: %s", job.Name)
		}
	}

	c.Start()

	time.Sleep(5 * time.Minute)

	c.Stop()
	return nil
}
