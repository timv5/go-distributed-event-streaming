package model

import "time"

type Message struct {
	MessageId string    `gorm:"type:bigint;primary_key" sql:"messageId"`
	Header    string    `gorm:"not null" sql:"header"`
	Body      string    `gorm:"not null" sql:"body"`
	CreatedAt time.Time `gorm:"not null" sql:"createdAt"`
	UpdatedAt time.Time `gorm:"not null" sql:"UpdatedAt"`
	Status    string    `gorm:"not null" sql:"status"`
}
