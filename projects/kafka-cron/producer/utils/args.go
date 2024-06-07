package utils

import "flag"

func Args() (string, string, string) {
	topic1 := flag.String("topic1", "topic-cluster-a", "Kafka topic to send messages to")
	topic2 := flag.String("topic2", "topic-cluster-b", "Kafka topic to send messages to")
	brokers := flag.String("brokers", "localhost:9092", "Kafka brokers to connect to")

	flag.Parse()

	return *topic1, *topic2, *brokers

}
