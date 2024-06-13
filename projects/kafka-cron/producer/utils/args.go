package utils

import "flag"

func Args() (string, string, string) {
	topics := flag.String("topics", "topic-cluster-a,topic-cluster-b", "List of Kafka topics separated by comma")
	clusters := flag.String("clusters", "cluster-a,cluster-b", "List of Kafka clusters separated by comma")
	brokers := flag.String("brokers", "localhost:9092", "Kafka brokers to connect to")
	flag.Parse()

	return *topics, *clusters, *brokers

}
