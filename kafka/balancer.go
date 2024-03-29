package kafka

import (
	"math/rand"

	"github.com/segmentio/kafka-go"
)

type randomBalancer struct{}

func (b randomBalancer) Balance(msg kafka.Message, partitions ...int) (partition int) {
	return partitions[rand.Int()%len(partitions)]
}

type PartitionBalancer struct {
	rr randomBalancer
}

func (balancer *PartitionBalancer) Balance(msg kafka.Message, partitions ...int) int {
	if msg.Partition == 0 || msg.Partition > len(partitions) {
		return balancer.rr.Balance(msg, partitions...)
	}
	return msg.Partition - 1
}
