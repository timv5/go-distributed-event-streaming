package main

import (
	"consumer/configs"
	"consumer/repository"
	"consumer/service"
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
	gormDB, err := configs.ConnectToDB(&config)
	if err != nil {
		panic("Failed to connect to DB")
	}

	// initialize repository
	messageRepository := repository.NewMessageRepository()
	messageHistoryRepository := repository.NewMessageHistoryRepository()

	// initialize service
	consumerService := service.NewConsumerService(&config, messageRepository, messageHistoryRepository, gormDB)

	initializeRMQ(&config, consumerService)
}

func initializeRMQ(config *configs.Config, consumerService *service.ConsumerService) {
	// set rmq
	conn, err := amqp.Dial(config.RMQUrl)
	if err != nil {
		panic("Could not initialize RMQ")
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			panic("Could not initialize RMQ")
		}
	}(conn)

	log.Println("Successfully connected to RMQ")

	// connect to channel
	ch, err := conn.Channel()
	if err != nil {
		panic("Cannot connect to RMQ channel")
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			panic("Cannot connect to RMQ channel")
		}
	}(ch)

	queue, err := ch.QueueDeclare(config.RMQQueueName, false, false, false, false, nil)
	if err != nil {
		panic("Cannot connect to Queue")
	}
	log.Println(queue)

	msg, err := ch.Consume(config.RMQQueueName, "", true, false, false, false, nil)

	forever := make(chan bool)
	go func() {
		for m := range msg {
			consumerService.HandleMessage(m)
		}
	}()
	<-forever
}
