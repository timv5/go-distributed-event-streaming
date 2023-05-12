package handler

import (
	"github.com/gin-gonic/gin"
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/dto/request"
	"go-distributed-event-streaming/service"
	"gorm.io/gorm"
	"net/http"
)

type MessageHandler struct {
	postgresDB     *gorm.DB
	messageService *service.MessageService
	config         *configs.Config
}

func NewMessageHandler(postgresDB *gorm.DB, messageService *service.MessageService, config *configs.Config) MessageHandler {
	return MessageHandler{
		postgresDB:     postgresDB,
		messageService: messageService,
		config:         config,
	}
}

func (messageHandler MessageHandler) SendMessage(ctx *gin.Context) {
	var messagePayload *request.MessageRequest
	if err := ctx.ShouldBindJSON(&messagePayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if messagePayload.Body == "" || messagePayload.Header == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "wrong request params"})
		return
	}

	messageResponse, err := messageHandler.messageService.SaveMessage(messagePayload.Header, messagePayload.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status:": "error", "message": err})
		return
	} else {
		ctx.JSON(http.StatusOK, messageResponse)
	}
}
