package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/awesome-sphere/as-booking/utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

var TOPIC string
var PARTITION int
var KAFKA_ADDR string

func PushMessage(message_value *WriterInterface) (bool, error) {
	config := kafka.WriterConfig{
		Brokers:          []string{KAFKA_ADDR},
		Topic:            TOPIC,
		Balancer:         &kafka.LeastBytes{},
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	writer_connector := kafka.NewWriter(config)
	defer writer_connector.Close()

	new_byte_buffer := new(bytes.Buffer)
	json.NewEncoder(new_byte_buffer).Encode(message_value)

	err := writer_connector.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(strconv.Itoa(message_value.TheaterId)),
			Value: new_byte_buffer.Bytes(),
		},
	)
	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	return true, nil
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
	conn, err := kafka.Dial("tcp", KAFKA_ADDR)
	if err != nil {
		panic(err.Error())
	}
	return conn
}

func InitKafkaTopic() {
	TOPIC = utils.GetenvOr("KAFKA_TOPIC", "ah")
	KAFKA_ADDR = utils.GetenvOr("KAFKA_ADDR", "localhost:9092")
	partition_num, err := strconv.Atoi(utils.GetenvOr("KAFKA_TOPIC_PARTITION", "5"))
	if err != nil {
		panic(err.Error())
	}
	PARTITION = partition_num

	conn := ConnectKafka()
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
}
