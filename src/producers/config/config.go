package config

import (
	"fmt"

	"main/src/goDotEnvVariable"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func Kafka() *kafka.Producer {
	brokers := goDotEnvVariable.GoDotEnvVariable("KAFKA_1") + ", " +
		goDotEnvVariable.GoDotEnvVariable("KAFKA_2") + ", " + goDotEnvVariable.GoDotEnvVariable("KAFKA_3")

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers, // producer can find the Kakfa cluster
		// "client.id": socket.gethostname(),	// easily correlate requests on the broker with the client instance
		"acks": "all", // acks=all;
		// "batch.size": 10000000,
		// "batch.num.messages": 100000,
		// "linger.ms": 180,
		// "compression.type": "lz4",
	})
	if err != nil {
		// When a connection error occurs,
		// a panic occurs and the system is shut down
		fmt.Printf("❗️Failed to create producer: %s", err)
		panic(err)
	}

	return p
}
