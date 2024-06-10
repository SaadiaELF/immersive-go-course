package models

type CronJob struct {
	Id            string `json:"id"`
	Schedule      string `json:"schedule"`
	Command       string `json:"command"`
	Cluster       string `json:"cluster"`
	Retries       int    `json:"retries"`
	RetryTopic    string `json:"retry_topic"`
	RetryInterval int    `json:"retry_interval"`
}
