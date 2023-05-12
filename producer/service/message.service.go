package service

import (
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/dto/response"
	"go-distributed-event-streaming/model"
	"go-distributed-event-streaming/producer"
	"go-distributed-event-streaming/repository"
)

type Message interface {
	SaveMessage(message *model.Message) (response *response.MessageResponse, err error)
}

type MessageService struct {
	conf              *configs.Config
	messageRepository *repository.MessageRepository
}

func NewMessageService(config *configs.Config, messageRepository *repository.MessageRepository) *MessageService {
	return &MessageService{conf: config, messageRepository: messageRepository}
}

func (mess *MessageService) SaveMessage(header string, body string) (response.MessageResponse, error) {
	// save message
	var savedMessage model.Message
	savedMessage, err := mess.messageRepository.SaveMessage(header, body)
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
		ID:        savedMessage.ID,
		CreatedAt: savedMessage.CreatedAt,
		Header:    savedMessage.Header,
		Body:      savedMessage.Body,
	}, nil
}
