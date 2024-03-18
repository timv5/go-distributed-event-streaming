package service

import (
	"consumer/configs"
	"consumer/repository"
	"consumer/rmq"
	"encoding/json"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
)

type ConsumerInterface interface {
	HandleMessage(msg amqp.Delivery)
}

type ConsumerService struct {
	conf                     *configs.Config
	postgresDB               *gorm.DB
	messageRepository        *repository.MessageRepository
	messageHistoryRepository *repository.MessageHistoryRepository
}

func NewConsumerService(config *configs.Config, messageRepository *repository.MessageRepository,
	messageHistoryRepository *repository.MessageHistoryRepository, postgresDB *gorm.DB) *ConsumerService {
	return &ConsumerService{
		conf:                     config,
		postgresDB:               postgresDB,
		messageRepository:        messageRepository,
		messageHistoryRepository: messageHistoryRepository,
	}
}

func (cons *ConsumerService) HandleMessage(msg amqp.Delivery) {
	response := &rmq.Message{}
	err := json.Unmarshal(msg.Body, response)
	if err != nil {
		log.Printf("ERROR: fail unmarshl: %s", msg.Body)
		return
	}

	tx := cons.postgresDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updatedMessage, err := cons.messageRepository.Update(tx, response)
	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: cannot update a message")
		return
	}

	_, err = cons.messageHistoryRepository.Save(tx, updatedMessage.MessageId, "SENT", "RECEIVED")
	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: cannot insert message history for message %s", updatedMessage.MessageId)
		return
	}

	tx.Commit()
	log.Printf("Successfully update: %s", updatedMessage)
}
