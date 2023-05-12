package rmq

import "time"

type Message struct {
	SentAt time.Time
	Header string
	Body   string
	ID     string
}
