package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9094"},
		Topic:   TOPIC,
		GroupID: "booking-consumer-1",
	})

	for {
		msg, err := r.ReadMessage(ctx)
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
