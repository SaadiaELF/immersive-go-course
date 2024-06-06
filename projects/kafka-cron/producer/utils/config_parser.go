package utils

import (
	"bufio"
	"fmt"
	"kafka-cron/pkg/models"
	"os"
	"strings"
)

func ParseConfig(cluster string) ([]models.CronJob, error) {

	path := "./config/crontab"
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var jobs []models.CronJob
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			return nil, fmt.Errorf("invalid crontab line")
		}

		job := models.CronJob{
			Schedule: strings.Join(fields[0:6], " "),
			Command:  strings.Join(fields[6:], " "),
			Cluster:  cluster,
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
