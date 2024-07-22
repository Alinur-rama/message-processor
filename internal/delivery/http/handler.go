package http

import (
	"context"
	"log"
	"message-processor/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	service domain.MessageService
}

func NewMessageHandler(service domain.MessageService) *MessageHandler {
	return &MessageHandler{service: service}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if err := h.service.CreateMessage(ctx, input.Content); err != nil {
		log.Printf("Error creating message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message created successfully"})
}

func (h *MessageHandler) GetStatistics(c *gin.Context) {
	ctx := context.Background()
	total, processed, err := h.service.GetStatistics(ctx)
	if err != nil {
		log.Printf("Error getting statistics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_messages":     total,
		"processed_messages": processed,
	})
}
