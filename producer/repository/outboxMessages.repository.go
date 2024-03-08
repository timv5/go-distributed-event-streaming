package repository

import (
	uuid "github.com/satori/go.uuid"
	"go-distributed-event-streaming/model"
	"gorm.io/gorm"
	"time"
)

type OutboxMessageRepositoryInterface interface {
	Save(tx *gorm.DB, body []byte) (model.OutboxMessage, error)
	Fetch(tx *gorm.DB) (model.OutboxMessage, error)
}

type OutboxMessageRepository struct{}

func NewOutboxMessageRepository() *OutboxMessageRepository {
	return &OutboxMessageRepository{}
}

func (repo *OutboxMessageRepository) Save(tx *gorm.DB, body []byte) (model.OutboxMessage, error) {
	nowTime := time.Now()
	outboxMessageEntity := model.OutboxMessage{OutboxMessageId: uuid.NewV4().String(), CreatedAt: nowTime, Body: body}

	savedOutboxMessageEntity := tx.Create(&outboxMessageEntity)
	if savedOutboxMessageEntity.Error != nil {
		return model.OutboxMessage{}, savedOutboxMessageEntity.Error
	} else {
		return outboxMessageEntity, nil
	}
}

func (repo *OutboxMessageRepository) Fetch(tx *gorm.DB) ([]model.OutboxMessage, error) {
	var messages []model.OutboxMessage
	err := tx.Limit(5).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	var messageIDs []string
	for _, msg := range messages {
		messageIDs = append(messageIDs, msg.OutboxMessageId)
	}

	if err := tx.Where("outbox_message_id IN (?)", messageIDs).Delete(&model.OutboxMessage{}).Error; err != nil {
		return nil, err
	}

	return messages, nil
}
