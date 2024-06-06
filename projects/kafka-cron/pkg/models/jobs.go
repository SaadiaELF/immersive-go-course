package models

type CronJob struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
	Cluster  string `json:"cluster"`
}
