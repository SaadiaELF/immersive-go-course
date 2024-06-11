package models

import "github.com/prometheus/client_golang/prometheus"

var (
	CronJobLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "cron_job_latency_seconds",
		Help:    "Latency of cron job execution in seconds.",
		Buckets: prometheus.DefBuckets,
	})
)
