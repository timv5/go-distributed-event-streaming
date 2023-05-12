package response

import "time"

type MessageResponse struct {
	ID        string
	CreatedAt time.Time
	Header    string
	Body      string
}
