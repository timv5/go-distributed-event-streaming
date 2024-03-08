package service

import (
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/dto/response"
	"go-distributed-event-streaming/repository"
	"go-distributed-event-streaming/utils"
	"gorm.io/gorm"
)

type Message interface {
	SaveMessage(header string, body string) (response.MessageResponse, error)
}

type MessageService struct {
	conf                     *configs.Config
	messageRepository        *repository.MessageRepository
	messageHistoryRepository *repository.MessageHistoryRepository
	outboxMessageRepository  *repository.OutboxMessageRepository
	postgresDB               *gorm.DB
}

func NewMessageService(config *configs.Config, messageRepository *repository.MessageRepository,
	messageHistoryRepository *repository.MessageHistoryRepository, postgresDB *gorm.DB, outboxMessageRepository *repository.OutboxMessageRepository) *MessageService {
	return &MessageService{conf: config, messageRepository: messageRepository, messageHistoryRepository: messageHistoryRepository, postgresDB: postgresDB, outboxMessageRepository: outboxMessageRepository}
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

	// save to outbox, so we can pick it up later
	savedMessageJson, err := utils.MarshalJSON(&savedMessage)
	_, err = mess.outboxMessageRepository.Save(tx, savedMessageJson)

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return response.MessageResponse{}, err
	}

	return response.MessageResponse{
		ID:        savedMessage.MessageId,
		CreatedAt: savedMessage.CreatedAt,
		Header:    savedMessage.Header,
		Body:      savedMessage.Body,
	}, nil
}
