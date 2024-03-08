package model

import "time"

type OutboxMessage struct {
	OutboxMessageId string    `gorm:"not null" sql:"outboxMessageId"`
	CreatedAt       time.Time `gorm:"not null" sql:"createdAt"`
	Body            []byte    `gorm:"not null" sql:"body"`
}
