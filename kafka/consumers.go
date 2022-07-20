package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/awesome-sphere/as-booking/db"
	"github.com/segmentio/kafka-go"
)

func initReader(topic string, groupID string, groupBalancers []kafka.GroupBalancer) *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers:        []string{KAFKA_ADDR},
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		GroupID:        groupID,
		GroupBalancers: groupBalancers,
	}
	r := kafka.NewReader(config)
	return r
}

func readFromReader(r *kafka.Reader) {
	defer r.Close()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Error while reading message: %v", err)
			continue
		}

		fmt.Printf("Reading message at topic [%v] partition [%v] offset [%v]: %s", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		db.UpdateStatus(msg.Topic, msg.Value)
	}
}

func Consume(topic string, groupID string, partition int) {
	groupBalancers := make([]kafka.GroupBalancer, 0)
	groupBalancers = append(groupBalancers, kafka.RangeGroupBalancer{})

	readers := make([]*kafka.Reader, 0)
	for i := 0; i < partition; i++ {
		readers = append(readers, initReader(topic, groupID, groupBalancers))
	}
	for _, reader := range readers {
		go readFromReader(reader)
	}
}
