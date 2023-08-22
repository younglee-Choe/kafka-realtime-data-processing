package main

import (
	"fmt"
	"main/src/consumer/config"
)

func main() {
	c := config.Kafka()

	// Subscribe to Kafka's topic to consume the decoded data
	c.SubscribeTopics([]string{"leele-topic-rekeyed"}, nil)
	defer c.Close()

	for {
		// ReadMessage polls the consumer for a message
		msg, err := c.ReadMessage(-1)
		if err == nil {
			// msg.TopicPartition provides partition-specific information (such as topic, partition and offset).
			fmt.Printf("✅ Received message %s: key=%s, vlaue=%s\n", msg.TopicPartition, msg.Key, msg.Value)
		} else {
			fmt.Printf("❗️ Consumer error: %v (%v)\n", err, msg)
		}
	}
}
