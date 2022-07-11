package kafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

var KafkaLeader *kafka.Conn

func InitKafkaTopic() {
	topic := "test3"
	partition := 5

	// Connect kafka broker
	conn, err := kafka.Dial("tcp", "localhost:9094")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	// Create Kafka topic
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partition,
			ReplicationFactor: 1,
		},
	}

	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}

	// List topic
	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}
	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}
	KafkaLeader = conn

}
