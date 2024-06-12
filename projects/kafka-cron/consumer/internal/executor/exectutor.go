package executor

import (
	"fmt"
	"kafka-cron/pkg/models"
	"os"
	"os/exec"
	"time"
)

func Execute(job models.CronJob) error {
	cmd := exec.Command("sh", "-c", job.Command)
	job.Latency = time.Since(job.StartTime)

	models.CronJobLatency.Observe(float64(job.Latency.Milliseconds()))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Executing job:%s\n", job.Id)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute job: %w", err)
	}
	fmt.Printf("Job executed successfully with latency: %s\n", job.Latency)
	return nil
}
