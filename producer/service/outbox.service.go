package service

import (
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/model"
	"go-distributed-event-streaming/producer"
	"go-distributed-event-streaming/utils"
	"log"
)

func HandleOutboxMessages(config *configs.Config, messages []model.OutboxMessage) error {
	for _, message := range messages {
		marshalledMessage, err := utils.UnmarshalJSON(message.Body)
		if err != nil {
			log.Printf("Error marshaling a message")
			return err
		} else {
			// publish message
			rmqProducer := producer.RMQProducer{
				ExchangeKey:      config.RMQExchangeKey,
				QueueName:        config.RMQQueueName,
				ConnectionString: config.RMQUrl,
			}

			rmqProducer.ProduceMessage(marshalledMessage)
		}
	}

	return nil
}
