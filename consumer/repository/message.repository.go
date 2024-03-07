package repository

import (
	"consumer/model"
	"consumer/rmq"
	"gorm.io/gorm"
	"time"
)

type MessageRepositoryInterface interface {
	Update(tx *gorm.DB, message *rmq.Message) (model.Message, error)
}

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (repo *MessageRepository) Update(tx *gorm.DB, message *rmq.Message) (model.Message, error) {
	savedMessage := tx.Model(&model.Message{}).Where("message_id = ?", message.ID).Updates(model.Message{UpdatedAt: time.Now(), Status: "RECEIVED"})

	if savedMessage.Error != nil {
		return model.Message{}, savedMessage.Error
	} else {
		return model.Message{MessageId: message.ID}, nil
	}
}
