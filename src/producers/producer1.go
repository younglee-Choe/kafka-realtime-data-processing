package producers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"main/src/generateMessage"
	"main/src/producers/config"
	"main/src/structures"
)

func Producer1() {
	// Producer configuration as a Kafka client
	p := config.Kafka()
	defer p.Close()

	deliveryChan := make(chan kafka.Event)
	var startTime time.Time

	for {
		// Go-routine to handle message delivery reports and
		// possibly other event types (errors, stats, etc)
		go func() {
			// Executed when data arrives in Kafka
			deliveryReport := <-deliveryChan
			m := deliveryReport.(*kafka.Message)

			if m.TopicPartition.Error != nil {
				fmt.Printf("❗️ Delivery failed: %v\n", m.TopicPartition.Error)
			} else if string(m.Value) == "Over rate limit" {
				fmt.Println("❗️ Over rate limit")
				time.Sleep(5 * time.Second)
			} else {
				// Check acks
				fmt.Printf("Delivered message to topic %s[%d] at offset %v with key '%s'\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset, m.Key)
			}

			endTime := time.Now()
			latency := endTime.Sub(startTime)
			fmt.Printf("latency: %s\n", latency)
		}()

		// Read JSON file to push to Kafka
		filePath := "./dummy_data/small_dummy.json"
		message := generateMessage.GenerateMessage(filePath)

		startTime = time.Now()

		var data structures.BlockData
		data.SourceData = message
		data.Length = len(message)

		// When sending data to kafka, marshal it and send it in byte array format
		value, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		topic := "leele-topic"
		key := "1producer"

		// Push message to Kafka with key by specifying topic
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(key),
			Value:          value,
		}, deliveryChan) // deliveryChan is for checking acks
		if err != nil {
			fmt.Printf("❗️ Failed to produce message: %s\n", err)
			panic(err)
		}

		// Wait for message deliveries before shutting down
		p.Flush(15 * 1000)
	}
}
