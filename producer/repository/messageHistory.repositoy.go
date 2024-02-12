package repository

import (
	"go-distributed-event-streaming/model"
	"gorm.io/gorm"
	"log"
	"time"
)

type MessageHistoryRepositoryInterface interface {
	Save(messageId string, fromStatus string, toStatus string) (model.MessageHistory, error)
}

type MessageHistoryRepository struct {
	postgresDB *gorm.DB
}

func NewMessageHistoryRepository(postgresDB *gorm.DB) *MessageHistoryRepository {
	return &MessageHistoryRepository{postgresDB: postgresDB}
}

func (repo *MessageHistoryRepository) Save(messageId string, fromStatus string, toStatus string) (model.MessageHistory, error) {
	nowTime := time.Now()
	messageHistoryEntity := model.MessageHistory{
		MessageId:  messageId,
		CreatedAt:  nowTime,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
	}

	savedMessageHistory := repo.postgresDB.Create(&messageHistoryEntity)
	if savedMessageHistory.Error != nil {
		return model.MessageHistory{}, savedMessageHistory.Error
	} else {
		log.Printf("Message history saved with id %s\n", messageHistoryEntity.MessageId)
		return messageHistoryEntity, nil
	}
}
