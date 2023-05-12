package request

type MessageRequest struct {
	Header string `json:"header" binding:"required"`
	Body   string `json:"body" binding:"required"`
}
