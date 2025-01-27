package main

import (
	"log"
	"os"

	DB "github.com/M0rfes/go-chat-ms/pkg/db"
	"github.com/M0rfes/go-chat-ms/pkg/db/models"
	ms "github.com/M0rfes/go-chat-ms/pkg/message-queue"
	"github.com/M0rfes/go-chat-ms/pkg/message-queue/consumers"
	msm "github.com/M0rfes/go-chat-ms/pkg/message-queue/models"
)

var (
	dbConnectionURL string
	kafkaURL        string
)

func init() {
	dbConnectionURL = os.Getenv("DB_URL")
	if dbConnectionURL == "" {
		panic("DB_CONNECTION_URL is required")
	}
	println("============>", dbConnectionURL)
	kafkaURL = os.Getenv("KAFKA_URL")
	if kafkaURL == "" {
		panic("KAFKA_URL is required")
	}

}

func main() {
	db := DB.NewDB(dbConnectionURL)
	err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	messageConsumer, err := consumers.NewMessage(kafkaURL, ms.MessagesTopic)
	if err != nil {
		panic(err)
	}

	msgChan := make(chan *msm.Message)
	ackChan := make(chan struct{})

	go messageConsumer.Consume(msgChan, ackChan)

	for msg := range msgChan {
		if err = db.CreteMessage(&models.Message{
			From:    msg.From,
			Content: msg.Message,
		}); err == nil {
			ackChan <- struct{}{}
		}
		log.Print(err.Error())
	}
}
