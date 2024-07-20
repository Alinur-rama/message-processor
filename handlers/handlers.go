package handlers

import (
	"encoding/json"
	"message-processor/database"
	"message-processor/kafka"
	"message-processor/models"
	"net/http"

	"github.com/IBM/sarama"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Сохранение сообщения в базу данных
	err = database.DB.QueryRow("INSERT INTO messages (content) VALUES ($1) RETURNING id", msg.Content).Scan(&msg.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправка сообщения в Kafka
	_, _, err = kafka.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: "messages",
		Value: sarama.StringEncoder(msg.Content),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		TotalMessages     int `json:"total_messages"`
		ProcessedMessages int `json:"processed_messages"`
	}

	err := database.DB.QueryRow("SELECT COUNT(*), COUNT(*) FILTER (WHERE processed = true) FROM messages").Scan(&stats.TotalMessages, &stats.ProcessedMessages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}
