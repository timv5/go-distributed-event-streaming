package service

import (
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/dto/response"
	"go-distributed-event-streaming/producer"
	"go-distributed-event-streaming/repository"
)

type Message interface {
	SaveMessage(header string, body string) (response.MessageResponse, error)
}

type MessageService struct {
	conf                     *configs.Config
	messageRepository        *repository.MessageRepository
	messageHistoryRepository *repository.MessageHistoryRepository
}

func NewMessageService(config *configs.Config, messageRepository *repository.MessageRepository, messageHistoryRepository *repository.MessageHistoryRepository) *MessageService {
	return &MessageService{conf: config, messageRepository: messageRepository, messageHistoryRepository: messageHistoryRepository}
}

func (mess *MessageService) SaveMessage(header string, body string) (response.MessageResponse, error) {
	// save message
	savedMessage, err := mess.messageRepository.SaveMessage(header, body)
	if err != nil {
		return response.MessageResponse{}, err
	}

	_, err = mess.messageHistoryRepository.Save(savedMessage.MessageId, "", "SENT")
	if err != nil {
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
