package repository

import (
	uuid "github.com/satori/go.uuid"
	"go-distributed-event-streaming/model"
	"gorm.io/gorm"
	"time"
)

type MessageRepositoryInterface interface {
	SaveMessage(tx *gorm.DB, header string, body string) (model.Message, error)
}

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (repo *MessageRepository) SaveMessage(tx *gorm.DB, header string, body string) (model.Message, error) {
	nowTime := time.Now()
	createMessage := model.Message{MessageId: uuid.NewV4().String(), CreatedAt: nowTime, Body: body, Header: header, Status: "SENT"}

	savedMessage := tx.Create(&createMessage)
	if savedMessage.Error != nil {
		return model.Message{}, savedMessage.Error
	} else {
		return createMessage, nil
	}
}
