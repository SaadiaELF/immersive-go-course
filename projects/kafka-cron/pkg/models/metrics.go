package models

import "github.com/prometheus/client_golang/prometheus"

var (
	CronJobLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "cron_job_latency_milliseconds",
		Help:    "Latency of cron job execution in milliseconds.",
		Buckets: prometheus.LinearBuckets(0, 10, 10),
	}, []string{"cluster"})
)

var (
	CronJobCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cron_job_count",
		Help: "Number of cron jobs executed.",
	}, []string{"cluster"})
)

var (
	CronJobErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cron_job_error_count",
		Help: "Number of cron jobs that failed.",
	}, []string{"cluster"})
)
