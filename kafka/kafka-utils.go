package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/awesome-sphere/as-booking/utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

var TOPIC string
var PARTITION int

func PushMessage(topic_name string, key_parition string, value *WriterInterface) {
	// writer_connector := &kafka.Writer{
	// 	Addr:     kafka.TCP("localhost:9094"),
	// 	Topic:    topic_name,
	// 	Balancer: &kafka.LeastBytes{},
	// }
	config := kafka.WriterConfig{
		Brokers:          []string{"localhost:9092"},
		Topic:            topic_name,
		Balancer:         &kafka.LeastBytes{},
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	writer_connector := kafka.NewWriter(config)
	defer writer_connector.Close()

	new_byte_buffer := new(bytes.Buffer)
	json.NewEncoder(new_byte_buffer).Encode(value)

	err := writer_connector.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(key_parition),
			Value: new_byte_buffer.Bytes(),
		},
	)
	if err != nil {
		panic(err.Error())
	}
}

func ListTopic(connector *kafka.Conn) map[string]*TopicInterface {
	partitions, err := connector.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}
	m := make(map[string]*TopicInterface)
	for _, p := range partitions {
		if _, ok := m[p.Topic]; ok {
			m[p.Topic].Partition += 1
		} else {
			m[p.Topic] = &TopicInterface{Partition: 1}
		}
	}
	return m

}

func isTopicExist(connector *kafka.Conn, topic_name string) bool {
	topic_map := ListTopic(connector)
	_, ok := topic_map[topic_name]
	return ok

}

func ConnectKafka() *kafka.Conn {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	return conn
}

func InitKafkaTopic() {
	TOPIC := utils.GetenvOr("KAFKA_TOPIC", "ah")
	PARTITION, err := strconv.Atoi(utils.GetenvOr("KAFKA_TOPIC_PARTITION", "5"))
	if err != nil {
		panic(err.Error())
	}

	conn := ConnectKafka()
	fmt.Printf("\n\n\nConnected\n\n\n")
	defer conn.Close()

	if !isTopicExist(conn, TOPIC) {
		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             TOPIC,
				NumPartitions:     PARTITION,
				ReplicationFactor: 1,
			},
		}

		err := conn.CreateTopics(topicConfigs...)
		if err != nil {
			panic(err.Error())
		}
	}
	test := &WriterInterface{
		UserID:     1,
		TimeSlotId: 1,
		TheaterId:  1,
		SeatNumber: 1,
		SeatTypeId: 1,
	}
	PushMessage(TOPIC, "1", test)
}
