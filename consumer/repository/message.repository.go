package repository

import (
	"consumer/model"
	"consumer/rmq"
	"gorm.io/gorm"
	"time"
)

type MessageRepositoryInterface interface {
	Update(message *rmq.Message) (model.Message, error)
}

type MessageRepository struct {
	postgresDB *gorm.DB
}

func NewMessageRepository(postgresDB *gorm.DB) *MessageRepository {
	return &MessageRepository{postgresDB: postgresDB}
}

func (repo *MessageRepository) Update(message *rmq.Message) (model.Message, error) {
	savedMessage := repo.postgresDB.Model(&model.Message{}).Where("message_id = ?", message.ID).Updates(model.Message{UpdatedAt: time.Now(), Status: "RECEIVED"})

	if savedMessage.Error != nil {
		return model.Message{}, savedMessage.Error
	} else {
		return model.Message{MessageId: message.ID}, nil
	}
}
