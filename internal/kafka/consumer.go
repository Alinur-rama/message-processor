package kafka

import (
	"context"
	"log"
	"message-processor/internal/domain"

	"github.com/IBM/sarama"
)

type Consumer interface {
	Start(ctx context.Context) error
}

type KafkaConsumer struct {
	consumer sarama.Consumer
	topic    string
	service  domain.MessageService
}

func NewConsumer(brokers []string, topic string, service domain.MessageService) (Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
		service:  service,
	}, nil
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	log.Println("Kafka consumer started")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			if err := c.processMessage(ctx, string(msg.Value)); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from Kafka consumer: %v", err)
		case <-ctx.Done():
			log.Println("Shutting down Kafka consumer")
			return nil
		}
	}
}

func (c *KafkaConsumer) processMessage(ctx context.Context, content string) error {
	log.Printf("Processing message: %s", content)
	return c.service.ProcessMessage(ctx, content)
}
