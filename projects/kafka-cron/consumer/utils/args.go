package utils

import "flag"

func Args() (string, string, string, bool) {
	topics := flag.String("topics", "topic-cluster-a,topic-cluster-b,topic3", "List of Kafka topics separated by comma")
	brokers := flag.String("brokers", "localhost:9092", "Kafka brokers to connect to")
	cluster := flag.String("cluster", "cluster-a", "Specify the cluster (cluster-a or cluster-b)")
	retry := flag.Bool("retry", false, "Retry sending messages to the topic")

	flag.Parse()

	return *topics, *brokers, *cluster, *retry
}
