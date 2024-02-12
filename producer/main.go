package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/handler"
	"go-distributed-event-streaming/repository"
	"go-distributed-event-streaming/route"
	"go-distributed-event-streaming/service"
	"log"
)

var (
	server                 *gin.Engine
	MessageController      handler.MessageHandler
	MessageRouteController route.MessageRouteHandler
)

func main() {
	// set configs
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic("Could not initialize app")
	}

	// connect to database
	configs.ConnectToDB(&config)

	// initialize repository
	messageRepository := repository.NewMessageRepository(configs.DB)
	messageHistoryRepository := repository.NewMessageHistoryRepository(configs.DB)

	// initialize service
	messageService := service.NewMessageService(&config, messageRepository, messageHistoryRepository)

	// initialize controllers and routes
	MessageController = handler.NewMessageHandler(configs.DB, messageService, &config)
	MessageRouteController = route.NewMessageRouteHandler(MessageController)

	server = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	MessageRouteController.MessageRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
