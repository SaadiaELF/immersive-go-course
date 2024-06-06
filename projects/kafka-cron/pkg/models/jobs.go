package models

type CronJob struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
	Cluster  string `json:"cluster"`
}
