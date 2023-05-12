package route

import (
	"github.com/gin-gonic/gin"
	"go-distributed-event-streaming/handler"
)

type MessageRouteHandler struct {
	messageHandler handler.MessageHandler
}

func NewMessageRouteHandler(messageHandler handler.MessageHandler) MessageRouteHandler {
	return MessageRouteHandler{
		messageHandler: messageHandler,
	}
}

func (h *MessageRouteHandler) MessageRoute(group *gin.RouterGroup) {
	router := group.Group("message")
	router.POST("/send", h.messageHandler.SendMessage)
}
