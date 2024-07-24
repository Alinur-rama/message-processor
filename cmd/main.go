package main

import (
	"context"
	"log"
	"message-processor/internal/database"
	"message-processor/internal/delivery/http"
	"message-processor/internal/kafka"
	"message-processor/internal/logger"
	"message-processor/internal/repository"
	"message-processor/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Init()

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	producer, err := kafka.NewProducer([]string{"kafka:9092"}, "messages")
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	repo := repository.NewPostgresMessageRepository(db)

	svc := service.NewMessageService(repo, producer)

	handler := http.NewMessageHandler(svc)

	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.POST("/messages", handler.CreateMessage)
		api.GET("/statistics", handler.GetStatistics)
	}

	consumer, err := kafka.NewConsumer([]string{"kafka:9092"}, "messages", svc)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	go consumer.Start(context.Background())

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
