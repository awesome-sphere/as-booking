package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/awesome-sphere/as-booking/kafka/writer_interface"
	"github.com/awesome-sphere/as-booking/utils"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

var BOOKING_TOPIC string
var CANCEL_TOPIC string
var BOOKING_PARTITION int
var CANCEL_PARTITION int
var KAFKA_ADDR string

type Result struct {
	IsCompleted bool
	Err         error
}

func Produce(message_value *writer_interface.BookingWriterInterface, result chan Result) {
	isCompleted, err := PushBookingMessage(message_value)
	result <- Result{IsCompleted: isCompleted, Err: err}
}

func PushBookingMessage(message_value *writer_interface.BookingWriterInterface) (bool, error) {
	config := kafka.WriterConfig{
		Brokers:          []string{KAFKA_ADDR},
		Topic:            BOOKING_TOPIC,
		Balancer:         &PartitionBalancer{},
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

func PushCancelMessage(message_value *writer_interface.CancelWriterInterface) (bool, error) {
	config := kafka.WriterConfig{
		Brokers:          []string{KAFKA_ADDR},
		Topic:            CANCEL_TOPIC,
		Balancer:         &PartitionBalancer{},
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

func doesTopicExist(connector *kafka.Conn, topic_name string) bool {
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
	BOOKING_TOPIC = utils.GetenvOr("KAFKA_BOOKING_TOPIC_NAME", "booking")
	CANCEL_TOPIC = utils.GetenvOr("KAFKA_CANCEL_TOPIC_NAME", "cancel_booking")
	KAFKA_ADDR = utils.GetenvOr("KAFKA_ADDR", "localhost:9092")
	partition_num, err := strconv.Atoi(utils.GetenvOr("KAFKA_BOOKING_TOPIC_PARTITION", "5"))
	if err != nil {
		panic(err.Error())
	}
	BOOKING_PARTITION = partition_num
	cancel_partition_num, err := strconv.Atoi(utils.GetenvOr("KAFKA_CANCEL_TOPIC_PARTITION", "5"))
	if err != nil {
		panic(err.Error())
	}
	CANCEL_PARTITION = cancel_partition_num

	conn := ConnectKafka()
	defer conn.Close()

	if !doesTopicExist(conn, BOOKING_TOPIC) {
		updateBookingtopicConfigs := []kafka.TopicConfig{
			{
				Topic:             BOOKING_TOPIC,
				NumPartitions:     BOOKING_PARTITION,
				ReplicationFactor: 1,
			},
		}
		err := conn.CreateTopics(updateBookingtopicConfigs...)
		if err != nil {
			panic("Booking Topic: " + err.Error())
		}
	}

	if !doesTopicExist(conn, CANCEL_TOPIC) {
		cancelTopicConfigs := []kafka.TopicConfig{
			{
				Topic:             CANCEL_TOPIC,
				NumPartitions:     CANCEL_PARTITION,
				ReplicationFactor: 1,
			},
		}
		err := conn.CreateTopics(cancelTopicConfigs...)
		if err != nil {
			panic("Cancel Topic: " + err.Error())
		}
	}
}
