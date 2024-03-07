package main

import (
	"consumer/configs"
	"consumer/repository"
	"consumer/rmq"
	"encoding/json"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
)

func main() {
	// set configs
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic("Could not initialize app")
	}

	// connect to database
	gormDB, err := configs.ConnectToDB(&config)
	if err != nil {
		panic("Failed to connect to DB")
	}

	// initialize repository
	messageRepository := repository.NewMessageRepository()
	messageHistoryRepository := repository.NewMessageHistoryRepository()

	initializeRMQ(&config, messageRepository, messageHistoryRepository, gormDB)
}

func initializeRMQ(config *configs.Config, messageRepository *repository.MessageRepository, messageHistoryRepository *repository.MessageHistoryRepository, postgresDB *gorm.DB) {
	// set rmq
	conn, err := amqp.Dial(config.RMQUrl)
	if err != nil {
		panic("Could not initialize RMQ")
	}
	defer conn.Close()

	log.Println("Successfully connected to RMQ")

	// connect to channel
	ch, err := conn.Channel()
	if err != nil {
		panic("Cannot connect to RMQ channel")
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(config.RMQQueueName, false, false, false, false, nil)
	if err != nil {
		panic("Cannot connect to Queue")
	}
	log.Println(queue)

	msg, err := ch.Consume(config.RMQQueueName, "", true, false, false, false, nil)

	forever := make(chan bool)
	go func() {
		for m := range msg {
			go handleMessage(m, messageRepository, messageHistoryRepository, postgresDB)
		}
	}()
	<-forever
}

func handleMessage(msg amqp.Delivery, messageRepository *repository.MessageRepository, messageHistoryRepository *repository.MessageHistoryRepository, postgresDB *gorm.DB) {
	response := &rmq.Message{}
	err := json.Unmarshal(msg.Body, response)
	if err != nil {
		log.Printf("ERROR: fail unmarshl: %s", msg.Body)
		return
	}

	tx := postgresDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updatedMessage, err := messageRepository.Update(tx, response)
	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: cannot update a message")
		return
	}

	_, err = messageHistoryRepository.Save(tx, updatedMessage.MessageId, "SENT", "RECEIVED")
	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: cannot insert message history for message %s", updatedMessage.MessageId)
		return
	}

	tx.Commit()
	log.Printf("Successfully update: %s", updatedMessage)
}
