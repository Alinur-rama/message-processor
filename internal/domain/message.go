package domain

import (
	"context"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Processed bool      `json:"processed"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageRepository interface {
	Create(ctx context.Context, message *Message) error
	MarkAsProcessed(ctx context.Context, content string) error
	GetStatistics(ctx context.Context) (int, int, error)
}

type MessageService interface {
	CreateMessage(ctx context.Context, content string) error
	ProcessMessage(ctx context.Context, content string) error
	GetStatistics(ctx context.Context) (int, int, error)
}
