package main

import (
	"context"
	"log"
	"message-processor/consumer"
	"message-processor/database"
	"message-processor/handlers"
	"message-processor/kafka"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	err = kafka.InitKafkaProducer()
	if err != nil {
		log.Fatal(err)
	}
	defer kafka.Producer.Close()

	r := mux.NewRouter()
	r.HandleFunc("/messages", handlers.CreateMessage).Methods("POST")
	r.HandleFunc("/statistics", handlers.GetStatistics).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем Kafka Consumer в отдельной горутине
	go consumer.StartConsumer(ctx)

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		log.Println("Server starting on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Ожидаем сигнала для грациозного завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Отменяем контекст для завершения Consumer
	cancel()

	// Создаем контекст с таймаутом для завершения сервера
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
