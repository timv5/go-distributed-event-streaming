package main

import (
	"consumer/configs"
	"consumer/repository"
	"consumer/rmq"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// set configs
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic("Could not initialize app")
	}

	// connect to database
	configs.ConnectToDB(&config)

	// initialize repository
	messageRepository := repository.NewMessageRepository(configs.DB)

	initializeRMQ(&config, messageRepository)
}

func initializeRMQ(config *configs.Config, messageRepository *repository.MessageRepository) {
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
			go handleMessage(m, messageRepository)
		}
	}()
	<-forever
}

func handleMessage(msg amqp.Delivery, messageRepository *repository.MessageRepository) {
	response := &rmq.Message{}
	err := json.Unmarshal(msg.Body, response)

	if err != nil {
		log.Printf("ERROR: fail unmarshl: %s", msg.Body)
		return
	}

	log.Printf("Successfully extracted: %s\n", response)
	updatedMessage, err := messageRepository.UpdateMessage(response)
	if err != nil {
		log.Printf("ERROR: cannot update a message")
	}
	log.Printf("Successfully update: %s", updatedMessage)
}
