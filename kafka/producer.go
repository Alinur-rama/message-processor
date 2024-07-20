package kafka

import (
	"log"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

func InitKafkaProducer() error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	var err error
	Producer, err = sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		return err
	}
	log.Println("Kafka producer initialized")
	return nil
}
