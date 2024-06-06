package utils

import "flag"

func Args() (string, string) {
	topic := flag.String("topic", "cron-jobs", "Kafka topic to send messages to")
	brokers := flag.String("brokers", "localhost:9092", "Kafka brokers to connect to")
	flag.Parse()

	return *topic, *brokers

}
