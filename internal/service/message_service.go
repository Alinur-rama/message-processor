package service

import (
	"context"
	"fmt"
	"message-processor/internal/domain"
	"message-processor/internal/kafka"
	"message-processor/internal/logger"

	"go.uber.org/zap"
)

type messageService struct {
	repo     domain.MessageRepository
	producer kafka.Producer
}

func NewMessageService(repo domain.MessageRepository, producer kafka.Producer) domain.MessageService {
	return &messageService{repo: repo, producer: producer}
}

func (s *messageService) CreateMessage(ctx context.Context, content string) error {
	logger.Log.Info("Creating new message", zap.String("content", content))

	message := &domain.Message{Content: content}
	if err := s.repo.Create(ctx, message); err != nil {
		logger.Log.Error("Failed to create message in repository", zap.Error(err))
		return fmt.Errorf("failed to create message: %w", err)
	}

	if err := s.producer.SendMessage(ctx, content); err != nil {
		logger.Log.Error("Failed to send message to Kafka", zap.Error(err))
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	logger.Log.Info("Message created successfully", zap.Int("id", message.ID))
	return nil
}

func (s *messageService) ProcessMessage(ctx context.Context, content string) error {
	return s.repo.MarkAsProcessed(ctx, content)
}

func (s *messageService) GetStatistics(ctx context.Context) (int, int, error) {
	return s.repo.GetStatistics(ctx)
}
