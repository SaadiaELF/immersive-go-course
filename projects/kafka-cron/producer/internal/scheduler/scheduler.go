package scheduler

import (
	"fmt"
	"kafka-cron/producer/pkg/models"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
)

func Scheduler(jobs []models.CronJob) {
	c := cron.New(cron.WithSeconds())

	// Schedule the jobs
	for _, job := range jobs {
		_, err := c.AddFunc(job.Schedule, func() {
			err := ExecuteJob(job)
			if err != nil {
				fmt.Println("error executing job:", job.Name)
			}
		})
		if err != nil {
			fmt.Println("error scheduling job:", job.Name)
		}
	}

	// Start the Cron job scheduler
	c.Start()

	// Wait for the Cron job to run
	time.Sleep(5 * time.Minute)

	// Stop the Cron job scheduler
	c.Stop()

}

func ExecuteJob(job models.CronJob) error {
	cmd := exec.Command("sh", "-c", job.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Executing job:", job.Name)
	return cmd.Run()
}
