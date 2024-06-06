package executor

import (
	"fmt"
	"kafka-cron/pkg/models"
	"os"
	"os/exec"
)

func Execute(job models.CronJob) error {
	cmd := exec.Command("sh", "-c", job.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Executing job:", job.Command)
	return cmd.Run()
}
