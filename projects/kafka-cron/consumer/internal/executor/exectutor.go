package executor

import (
	"fmt"
	"kafka-cron/pkg/models"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func Execute(job models.CronJob) error {
	n := rand.Intn(3)

	cmd := exec.Command("sh", "-c", job.Command)
	job.Latency = time.Since(job.StartTime)

	models.CronJobLatency.WithLabelValues(job.Cluster).Observe(float64(job.Latency.Milliseconds()))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Executing job:%s\n", job.Id)
	if n == 2 || n == 1{
		return fmt.Errorf("failed to execute job")
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute job: %w", err)
	}
	// fmt.Printf("Job executed successfully with latency: %s\n", job.Latency)
	return nil
}
