package consumer

import (
	"context"
	"log"
	"message-processor/database"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

func StartConsumer(ctx context.Context) {
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(kafkaBrokers, config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()

	topic := "messages"
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	log.Println("Kafka consumer started")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			processMessage(msg)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)
		case <-ctx.Done():
			log.Println("Shutting down consumer")
			return
		}
	}
}

func processMessage(msg *sarama.ConsumerMessage) {
	log.Printf("Received message: %s", string(msg.Value))

	// Обновляем статус сообщения в базе данных
	_, err := database.DB.Exec("UPDATE messages SET processed = true WHERE content = $1", string(msg.Value))
	if err != nil {
		log.Printf("Error updating message status: %v", err)
	} else {
		log.Println("Message marked as processed")
	}
}
