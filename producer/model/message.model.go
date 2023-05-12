package model

import "time"

type Message struct {
	ID        string    `gorm:"type:bigint;primary_key" sql:"id"`
	Header    string    `gorm:"not null" sql:"header"`
	Body      string    `gorm:"not null" sql:"body"`
	CreatedAt time.Time `gorm:"not null" sql:"createdAt"`
	UpdatedAt time.Time `sql:"updatedAt"`
	Status    string    `gorm:"not null" sql:"status"`
}
