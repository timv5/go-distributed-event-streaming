package producer

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go-distributed-event-streaming/dto/rmq"
	"go-distributed-event-streaming/model"
	"log"
	"time"
)

type RMQProducer struct {
	QueueName        string
	ConnectionString string
	ExchangeKey      string
}

func (pr RMQProducer) ProduceMessage(savedMessage *model.Message) {
	// set rmq
	conn, err := amqp.Dial(pr.ConnectionString)
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

	queue, err := ch.QueueDeclare(pr.QueueName, false, false, false, false, nil)
	if err != nil {
		panic("Cannot connect to Queue")
	}
	log.Println(queue)

	// prepare data
	message := rmq.Message{SentAt: time.Now(), Header: savedMessage.Header, Body: savedMessage.Body, ID: savedMessage.ID}
	stringMessage, err := json.Marshal(message)
	if err != nil {
		panic("Cannot marshal message to string")
	}

	err = ch.Publish(
		"",
		pr.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(string(stringMessage)),
		},
	)
	if err != nil {
		panic("Cannot publish")
	}

	log.Println("Successfully published message")
}
