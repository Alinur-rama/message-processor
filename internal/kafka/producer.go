package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(ctx context.Context, content string) error
	Close() error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string, topic string) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *KafkaProducer) SendMessage(ctx context.Context, content string) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(content),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	log.Printf("Message sent to Kafka: %s", content)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
