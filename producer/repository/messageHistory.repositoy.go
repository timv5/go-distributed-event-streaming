package repository

import (
	"go-distributed-event-streaming/model"
	"gorm.io/gorm"
	"log"
	"time"
)

type MessageHistoryRepositoryInterface interface {
	Save(tx *gorm.DB, messageId string, fromStatus string, toStatus string) (model.MessageHistory, error)
}

type MessageHistoryRepository struct{}

func NewMessageHistoryRepository() *MessageHistoryRepository {
	return &MessageHistoryRepository{}
}

func (repo *MessageHistoryRepository) Save(tx *gorm.DB, messageId string, fromStatus string, toStatus string) (model.MessageHistory, error) {
	nowTime := time.Now()
	messageHistoryEntity := model.MessageHistory{
		MessageId:  messageId,
		CreatedAt:  nowTime,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
	}

	savedMessageHistory := tx.Create(&messageHistoryEntity)
	if savedMessageHistory.Error != nil {
		return model.MessageHistory{}, savedMessageHistory.Error
	} else {
		log.Printf("Message history saved with id %s\n", messageHistoryEntity.MessageId)
		return messageHistoryEntity, nil
	}
}
