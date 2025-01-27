package producers

import (
	"encoding/json"

	topic "github.com/M0rfes/go-chat-ms/pkg/message-queue"
	model "github.com/M0rfes/go-chat-ms/pkg/message-queue/models"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Message interface {
	Publish(model.Message) error
	Close()
}

type message struct {
	producer *kafka.Producer
}

func NewMessage(url string) (Message, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": url})
	if err != nil {
		return nil, err
	}
	return &message{
		producer: producer,
	}, nil
}

func (m *message) Publish(msg model.Message) error {
	ms, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return m.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic.MessagesTopic, Partition: kafka.PartitionAny},
		Value:          ms,
	}, nil)
}

func (m *message) Close() {
	m.producer.Close()
}
