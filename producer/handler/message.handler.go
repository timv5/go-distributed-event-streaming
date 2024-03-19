package handler

import (
	"github.com/gin-gonic/gin"
	"go-distributed-event-streaming/client"
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/dto/request"
	"go-distributed-event-streaming/service"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type MessageHandlerInterface interface {
	SendMessage(ctx *gin.Context)
}

type MessageHandler struct {
	postgresDB     *gorm.DB
	messageService *service.MessageService
	config         *configs.Config
	restClient     *client.RestClient
}

func NewMessageHandler(postgresDB *gorm.DB, messageService *service.MessageService, config *configs.Config,
	restClient *client.RestClient) MessageHandler {
	return MessageHandler{
		postgresDB:     postgresDB,
		messageService: messageService,
		config:         config,
		restClient:     restClient,
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

	// call a rest client in a new go routine
	responseChannel := make(chan string)
	go messageHandler.restClient.CallService(responseChannel)

	messageResponse, err := messageHandler.messageService.SaveMessage(messagePayload.Header, messagePayload.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status:": "error", "message": err})
		return
	}

	// before completion get response from a rest client
	response := <-responseChannel
	log.Printf("Successfully got response from a rest client: %s", response)

	// return response
	ctx.JSON(http.StatusOK, messageResponse)
}
