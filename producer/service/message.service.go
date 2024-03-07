package service

import (
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/dto/response"
	"go-distributed-event-streaming/producer"
	"go-distributed-event-streaming/repository"
	"gorm.io/gorm"
)

type Message interface {
	SaveMessage(header string, body string) (response.MessageResponse, error)
}

type MessageService struct {
	conf                     *configs.Config
	messageRepository        *repository.MessageRepository
	messageHistoryRepository *repository.MessageHistoryRepository
	postgresDB               *gorm.DB
}

func NewMessageService(config *configs.Config, messageRepository *repository.MessageRepository,
	messageHistoryRepository *repository.MessageHistoryRepository, postgresDB *gorm.DB) *MessageService {
	return &MessageService{conf: config, messageRepository: messageRepository, messageHistoryRepository: messageHistoryRepository, postgresDB: postgresDB}
}

func (mess *MessageService) SaveMessage(header string, body string) (response.MessageResponse, error) {
	tx := mess.postgresDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// save message
	savedMessage, err := mess.messageRepository.SaveMessage(tx, header, body)
	if err != nil {
		tx.Rollback()
		return response.MessageResponse{}, err
	}

	_, err = mess.messageHistoryRepository.Save(tx, savedMessage.MessageId, "", "SENT")
	if err != nil {
		tx.Rollback()
		return response.MessageResponse{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return response.MessageResponse{}, err
	}

	// publish message
	rmqProducer := producer.RMQProducer{
		ExchangeKey:      mess.conf.RMQExchangeKey,
		QueueName:        mess.conf.RMQQueueName,
		ConnectionString: mess.conf.RMQUrl,
	}

	rmqProducer.ProduceMessage(&savedMessage)

	return response.MessageResponse{
		ID:        savedMessage.MessageId,
		CreatedAt: savedMessage.CreatedAt,
		Header:    savedMessage.Header,
		Body:      savedMessage.Body,
	}, nil
}
