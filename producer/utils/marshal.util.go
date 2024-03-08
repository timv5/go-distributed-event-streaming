package utils

import (
	"encoding/json"
	"go-distributed-event-streaming/model"
	"time"
)

func MarshalJSON(message *model.Message) ([]byte, error) {
	return json.Marshal(&struct {
		MessageID string    `json:"messageId"`
		Header    string    `json:"header"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Status    string    `json:"status"`
	}{
		MessageID: message.MessageId,
		Header:    message.Header,
		Body:      message.Body,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
		Status:    message.Status,
	})
}

func UnmarshalJSON(data []byte) (*model.Message, error) {
	var v struct {
		MessageID string    `json:"messageId"`
		Header    string    `json:"header"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Status    string    `json:"status"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &model.Message{
		MessageId: v.MessageID,
		Header:    v.Header,
		Body:      v.Body,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
		Status:    v.Status,
	}, nil
}
