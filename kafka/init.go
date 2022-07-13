package kafka

import (
	"strconv"

	"github.com/awesome-sphere/as-booking/utils"
	"github.com/segmentio/kafka-go"
)

var TOPIC string
var PARTITION int

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
	conn, err := kafka.Dial("tcp", "localhost:9094")
	if err != nil {
		panic(err.Error())
	}
	return conn
}

func InitKafkaTopic() {
	TOPIC := utils.GetenvOr("KAFKA_TOPIC", "my-awesome-topic")
	PARTITION, err := strconv.Atoi(utils.GetenvOr("KAFKA_TOPIC_PARTITION", "5"))
	if err != nil {
		panic(err.Error())
	}

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
