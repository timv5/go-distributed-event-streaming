package model

import "time"

type MessageHistory struct {
	MessageId  string    `gorm:"not null" sql:"messageId"`
	CreatedAt  time.Time `gorm:"not null" sql:"createdAt"`
	FromStatus string    `gorm:"not null" sql:"fromStatus"`
	ToStatus   string    `gorm:"not null" sql:"toStatus"`
}
