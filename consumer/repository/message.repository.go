package repository

import (
	"consumer/model"
	"consumer/rmq"
	"gorm.io/gorm"
	"time"
)

type MessageRepositoryInterface interface {
	SaveMessage(header string, body string)
}

type MessageRepository struct {
	postgresDB *gorm.DB
}

func NewMessageRepository(postgresDB *gorm.DB) *MessageRepository {
	return &MessageRepository{postgresDB: postgresDB}
}

func (repo *MessageRepository) UpdateMessage(message *rmq.Message) (model.Message, error) {
	savedMessage := repo.postgresDB.Model(&model.Message{}).Where("id = ?", message.ID).Updates(model.Message{UpdatedAt: time.Now(), Status: "RECEIVED"})

	//savedMessage := repo.postgresDB.Model(&model.Message{}).Where("id = ?", message.ID).Update("status", "RECEIVED")
	if savedMessage.Error != nil {
		return model.Message{}, savedMessage.Error
	} else {
		return model.Message{ID: message.ID}, nil
	}
}
