package config

import (
	"main/src/goDotEnvVariable"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func Kafka() *kafka.Consumer {
	brokers := goDotEnvVariable.GoDotEnvVariable("KAFKA_1") + ", " +
		goDotEnvVariable.GoDotEnvVariable("KAFKA_2") + ", " + goDotEnvVariable.GoDotEnvVariable("KAFKA_3")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		// When a connection error occurs, a panic occurs and the system is shut down
		panic(err)
	}

	return c
}
