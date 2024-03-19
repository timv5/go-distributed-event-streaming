package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-distributed-event-streaming/client"
	"go-distributed-event-streaming/configs"
	"go-distributed-event-streaming/handler"
	"go-distributed-event-streaming/producer"
	"go-distributed-event-streaming/repository"
	"go-distributed-event-streaming/route"
	"go-distributed-event-streaming/service"
	"log"
	"time"
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
	gormDB, err := configs.ConnectToDB(&config)
	if err != nil {
		panic("Failed to connect to DB")
	}

	// initialize repository
	messageRepository := repository.NewMessageRepository()
	messageHistoryRepository := repository.NewMessageHistoryRepository()
	outboxMessageRepository := repository.NewOutboxMessageRepository()

	// initialize rm handler
	rmqProducer := producer.NewRMQProducer(&config)

	// initialize service
	outboxService := service.NewOutboxService(rmqProducer)
	messageService := service.NewMessageService(&config, messageRepository, messageHistoryRepository, gormDB, outboxMessageRepository)
	restClient := client.NewRestClient()

	// scheduler
	go func() {
		for {
			tx := gormDB.Begin()
			messages, err := outboxMessageRepository.Fetch(tx)
			if err != nil {
				tx.Rollback()
				log.Fatalf("Failed to fetch messages: %v", err)
			}

			err = outboxService.HandleOutboxMessages(messages)
			if err != nil {
				tx.Rollback()
			}
			tx.Commit()

			time.Sleep(5 * time.Second)
		}
	}()

	// initialize controllers and routes
	MessageController = handler.NewMessageHandler(configs.DB, messageService, &config, &restClient)
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
