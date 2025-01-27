package consumers

import (
	"encoding/json"

	topic "github.com/M0rfes/go-chat-ms/pkg/message-queue"
	model "github.com/M0rfes/go-chat-ms/pkg/message-queue/models"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Message interface {
	Consume(Consumer[model.Message])
}

type message struct {
	url      string
	groupID  string
	consumer *kafka.Consumer
}

func NewMessage(url, groupID string) (Message, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": url,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	return &message{
		url:      url,
		groupID:  groupID,
		consumer: consumer,
	}, nil
}

func (m *message) Consume(consumer Consumer[model.Message]) {
	m.consumer.SubscribeTopics([]string{topic.MessagesTopic}, nil)
	defer m.consumer.Close()
	for {
		msg, err := m.consumer.ReadMessage(-1)
		if err == nil {
			consumer(nil, err)
			continue
		}
		message := &model.Message{}
		err = json.Unmarshal(msg.Value, message)
		if err == nil {
			consumer(nil, err)
		}
		consumer(message, nil)
	}
}
