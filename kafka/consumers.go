package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func Consume(topic_name string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{KAFKA_ADDR},
		Topic:   topic_name,
		GroupID: "booking-consumer",
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("could not read message: " + err.Error())
			break
		}
		fmt.Printf("message at topic [%v] partition [%v] offset [%v]: %s - %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader: ", err)
	}
}
