package models

import "github.com/prometheus/client_golang/prometheus"

var (
	CronJobLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "cron_job_latency_milliseconds",
		Help:    "Latency of cron job execution in milliseconds.",
		Buckets: prometheus.LinearBuckets(0, 5, 20),
	}, []string{"cluster"})
)
