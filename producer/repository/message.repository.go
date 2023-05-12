package repository

import (
	uuid "github.com/satori/go.uuid"
	"go-distributed-event-streaming/model"
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

func (repo *MessageRepository) SaveMessage(header string, body string) (model.Message, error) {
	nowTime := time.Now()
	createMessage := model.Message{ID: uuid.NewV4().String(), CreatedAt: nowTime, Body: body, Header: header, Status: "SENT"}

	savedMessage := repo.postgresDB.Create(&createMessage)
	if savedMessage.Error != nil {
		return model.Message{}, savedMessage.Error
	} else {
		return createMessage, nil
	}
}
