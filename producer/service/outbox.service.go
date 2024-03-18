package service

import (
	"go-distributed-event-streaming/model"
	"go-distributed-event-streaming/producer"
	"go-distributed-event-streaming/utils"
	"log"
)

type OutboxServiceInterface interface {
	HandleOutboxMessages(messages []model.OutboxMessage) error
}

type OutboxService struct {
	rmqProducer *producer.RMQProducer
}

func NewOutboxService(rmqProducer *producer.RMQProducer) *OutboxService {
	return &OutboxService{rmqProducer: rmqProducer}
}

func (om OutboxService) HandleOutboxMessages(messages []model.OutboxMessage) error {
	for _, message := range messages {
		marshalledMessage, err := utils.UnmarshalJSON(message.Body)
		if err != nil {
			log.Printf("Error marshaling a message")
			return err
		} else {
			om.rmqProducer.ProduceMessage(marshalledMessage)
		}
	}

	return nil
}
